package main

import (
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
	"net/http"
	"suno-api/common"
	"time"
)

func InitService() {
	common.InitTemplate()
	initTlsHTTPClient()

	if common.SessionID == "" {
		panic("suno account session_id not set")
	}
	if common.COOKIE == "" {
		panic("suno account cookie not set")
	}

	AccountInst = &Account{
		Certificate: SunoCert{
			SessionID: common.SessionID,
			Cookie:    common.COOKIE,
		},
	}

	startTaskWorker()

	common.SafeGoroutine(func() {
		for {
			time.Sleep(5 * time.Second)
			err := startAllKeepAlive()
			if err != nil {
				common.Logger.Error("Suno Keep-alive failed: " + err.Error())
			}
		}
	})
	common.SafeGoroutine(func() { // recover task
		recoverTasks()
	})
}

var TlsHTTPClient tls_client.HttpClient
var HTTPClient *http.Client

func initTlsHTTPClient() {
	jar := tls_client.NewCookieJar()
	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(360),
		tls_client.WithClientProfile(profiles.Chrome_120),
		tls_client.WithNotFollowRedirects(),
		tls_client.WithCookieJar(jar), // create cookieJar instance and pass it as argument
	}
	client, _ := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	TlsHTTPClient = client

	if common.Proxy != "" {
		err := TlsHTTPClient.SetProxy(common.Proxy)
		if err != nil {
			panic(err)
		}
	}
	HTTPClient = &http.Client{}
}
