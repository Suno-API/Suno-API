package main

import "C"
import (
	"encoding/json"
	"fmt"
	http "github.com/bogdanfinn/fhttp"
	"io"
	"strings"
	"suno-api/common"
	"time"
)

type Account struct {
	Certificate SunoCert `json:"certificate"`
	Msg         string   `json:"msg"`
}

type SunoCert struct {
	SessionID    string `json:"session_id"`
	Cookie       string `json:"cookie"`
	JWT          string `json:"jwt"`
	LastUpdate   int64  `json:"last_update"` // 最后更新时间，小于5秒，可以直接使用
	CreditsLeft  int    `json:"credits_left"`
	MonthlyLimit int    `json:"monthly_limit"`
	MonthlyUsage int    `json:"monthly_usage"`
	Period       string `json:"period"`
	IsActive     bool   `json:"is_active"`
}

var AccountInst *Account

func startAllKeepAlive() error {
	common.SafeGoroutine(func() {
		if err := accountKeepAlive(AccountInst); err != nil {
			common.Logger.Errorf("Suno Keep-alive 失败， err: %s", err.Error())
		}
	})
	return nil
}

func accountKeepAlive(acc *Account) error {
	// 保持账号活跃
	err := updateToken(acc)
	if err != nil {
		return err
	}
	return nil
}

const (
	exchangeTokenURL = "https://clerk.suno.com/v1/client/sessions/%s/tokens?_clerk_js_version=4.72.0-snapshot.vc141245"
)

func updateToken(a *Account) error {
	if a.Certificate.Cookie == "" {
		return fmt.Errorf("cookie 不可为空")
	}
	if a.Certificate.SessionID == "" {
		return fmt.Errorf("session_id 不可为空")
	}
	exchangeURL := fmt.Sprintf(exchangeTokenURL, a.Certificate.SessionID)
	req, err := http.NewRequest("POST", exchangeURL, nil)
	if err != nil {
		return err
	}

	for k, v := range CommonHeaders {
		req.Header.Set(k, v)
	}

	req.Header.Set("cookie", a.Certificate.Cookie)

	resp, err := TlsHTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			common.Logger.Errorw("body close", "err", err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	//var exchangeToken LastActiveToken
	exchangeToken := make(map[string]any)
	err = json.Unmarshal(body, &exchangeToken)
	if err != nil {
		return err
	}
	a.Certificate.JWT = common.Any2String(exchangeToken["jwt"])

	//get old cookies
	cookies := make(map[string]string)
	for _, cookie := range strings.Split(a.Certificate.Cookie, ";") {
		kv := strings.Split(cookie, "=")
		if len(kv) == 2 {
			cookies[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		}
	}

	setCookies := resp.Header.Values("Set-Cookie")
	for _, cookie := range setCookies {
		names := strings.Split(cookie, "; ")
		kv := strings.Split(names[0], "=")
		if len(kv) == 2 {
			cookies[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		}
	}
	var cookieArr []string
	for k, v := range cookies {
		cookieArr = append(cookieArr, k+"="+v)
	}
	a.Certificate.Cookie = strings.Join(cookieArr, "; ")

	err = getCredits(a)
	if err != nil {
		common.Logger.Errorf("GetCredits failed, err: %s", err.Error())
	}

	a.Certificate.LastUpdate = time.Now().Unix()

	return nil
}

func getCredits(account *Account) (err error) {
	resp, err := doRequest("GET", common.BaseUrl+"/api/billing/info/", nil, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + account.Certificate.JWT,
	})
	if err != nil {
		return
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			common.Logger.Errorw("body close", "err", err)
		}
	}(resp.Body)
	responseBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		err = fmt.Errorf("GetCredits failed, status: %d, body: %s", resp.StatusCode, string(responseBody))
		return
	}
	common.Logger.Debugw("GetCredits", "responseBody", string(responseBody))

	sunoResp := make(map[string]any)
	err = json.Unmarshal(responseBody, &sunoResp)
	if err != nil {
		return
	}

	account.Certificate.CreditsLeft = common.Any2Int(sunoResp["credits_left"])
	account.Certificate.MonthlyLimit = common.Any2Int(sunoResp["monthly_limit"])
	account.Certificate.MonthlyUsage = common.Any2Int(sunoResp["monthly_usage"])
	account.Certificate.Period = common.Any2String(sunoResp["period"])
	account.Certificate.IsActive = common.Any2Bool(sunoResp["is_active"])

	return nil
}
