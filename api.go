package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	http "github.com/bogdanfinn/fhttp"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"io"
	"strings"
	"suno-api/common"
	"suno-api/entity/po"
	"suno-api/lib/ginplus"
	"time"
)

var CommonHeaders = map[string]string{
	"Content-Type": "text/plain;charset=UTF-8",
	"User-Agent":   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36",
	"Referer":      "https://suno.com/",
	"Origin":       "https://suno.com",
	"Accept":       "*/*",
}

func doRequest(method, url string, body io.Reader, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	defer func() {
		if body != nil {
			err = req.Body.Close()
			if err != nil {
				common.Logger.Errorw("body close", "err", err)
			}
		}
	}()
	for k, v := range CommonHeaders {
		req.Header.Set(k, v)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	//req.Header.Set("Authorization", "Bearer "+po.SunoCert.JWT)

	resp, err := TlsHTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

var handles = map[string]func(reqBody []byte) (taskID string, relayErr *common.RelayError){
	po.TaskActionMusic:  submitGenSong, //该接口为自定义创建模式，灵活性高 // 自定义模式关闭的情况,通过对音乐的描述来生成歌曲。
	po.TaskActionLyrics: submitGenLyrics,
}

// @Summary Submit Suno song task
// @Schemes
// @Description
// @Tags suno
// @Accept json
// @Produce json
// @Param body body SubmitGenSongReq true "sumbmit generate song"
// @Success 200 {object} ginplus.DataResult{data=string} "task_id"
// @Router /api/submit/music [post]
func SubmitGenSong(c *gin.Context) {
	Submit(c)
}

// @Summary Submit Suno lyrics task
// @Schemes
// @Description
// @Tags suno
// @Accept json
// @Produce json
// @Param body body SubmitGenLyricsReq true "sumbmit generate lyrics"
// @Success 200 {object} ginplus.DataResult{data=string} "task_id"
// @Router /api/submit/lyrics [post]
func SubmitGenLyrics(c *gin.Context) {
	Submit(c)
}

func Submit(c *gin.Context) {
	action := c.Param("action")
	action = strings.ToUpper(action)
	c.Set("action", action)
	// choose action
	if handle, ok := handles[action]; ok {
		reqBody, err := io.ReadAll(c.Request.Body)
		if err != nil {
			common.ReturnErr(c, err, common.ErrCodeInvalidRequest, 400)
			return
		}
		taskID, relayErr := handle(reqBody)
		if relayErr != nil {
			common.ReturnErr(c, relayErr.Err, relayErr.Code, relayErr.StatusCode)
			return
		}
		c.JSON(200, ginplus.ApiRetSucc(taskID))
	} else {
		common.ReturnErr(c, fmt.Errorf("invalid action"), common.ErrCodeInvalidRequest, 400)
	}
}

// 自定义创建模式
func submitGenSong(reqBody []byte) (taskID string, relayErr *common.RelayError) {
	var err error
	var params SubmitGenSongReq
	err = json.Unmarshal(reqBody, &params)
	if err != nil {
		relayErr = common.WrapperErr(err, common.ErrCodeInvalidRequest, 400)
		return
	}
	common.Logger.Debugw("generateSong get", "requestBody", string(reqBody))

	reqData := make(map[string]interface{})
	reqData["make_instrumental"] = params.MakeInstrumental
	if params.Mv != "" {
		reqData["mv"] = params.Mv
	} else {
		reqData["mv"] = "chirp-v3-0"
	}
	if params.GptDescriptionPrompt != "" {
		reqData["gpt_description_prompt"] = params.GptDescriptionPrompt
		reqData["prompt"] = ""
	} else {
		reqData["prompt"] = params.Prompt
		reqData["title"] = params.Title
		reqData["tags"] = params.Tags
		reqData["continue_at"] = params.ContinueAt
		reqData["continue_clip_id"] = params.ContinueClipId
	}
	if params.ContinueClipId != nil && *params.ContinueClipId != "" { // 续写
		if params.TaskID == "" {
			relayErr = common.WrapperErr(fmt.Errorf("task_id is empty"), common.ErrCodeInvalidRequest, 400)
			return
		}
	}

	common.Logger.Debugw("generateSong do", "requestBody", reqData)

	requestBody, _ := json.Marshal(reqData)

	resp, err := doRequest("POST", common.BaseUrl+"/api/generate/v2/", bytes.NewReader(requestBody), map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + AccountInst.Certificate.JWT,
	})
	if err != nil {
		relayErr = common.WrapperErr(err, common.ErrCodeInternalError, 500)
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
		common.Logger.Errorw("generateSong", "status", resp.StatusCode, "requestBody", string(requestBody), "responseBody", string(responseBody))
		relayErr = common.WrapperErr(fmt.Errorf(resp.Status), common.ErrCodeInvalidRequest, resp.StatusCode)
		return
	}
	common.Logger.Debugw("generateSong", "responseBody", string(responseBody))

	var sunoResp GenSongResponse
	err = json.Unmarshal(responseBody, &sunoResp)
	if err != nil {
		relayErr = common.WrapperErr(err, common.ErrCodeInternalError, 500)
		return
	}
	// 保存任务

	if sunoResp.Status != "complete" {
		// 失败
		common.Logger.Errorw("generateSong failed", "responseBody", string(responseBody))
		err = fmt.Errorf("task failed")
		relayErr = common.WrapperErr(err, common.ErrCodeInternalError, 500)
		return
	}
	var tasks []po.TaskWithData[po.SunoSongs]
	task := po.TaskWithData[po.SunoSongs]{
		Action: po.TaskActionMusic,
		TaskID: common.GetUUID(),
	}
	task.Data = sunoResp.Clips
	tasks = append(tasks, task)
	// 保存任务
	err = po.AddTasks(tasks)
	if err != nil {
		relayErr = common.WrapperErr(err, common.ErrCodeInternalError, 500)
		return
	}
	AddTask(Task{
		Action: po.TaskActionMusic,
		ID:     task.TaskID,
	})
	// writer response
	//c.JSON(200, ginplus.ApiRetSucc(task.TaskID))
	return task.TaskID, nil
}

// 自定义模式关闭的情况
func submitGenLyrics(reqBody []byte) (taskID string, relayErr *common.RelayError) {
	var params SubmitGenLyricsReq
	err := json.Unmarshal(reqBody, &params)
	if err != nil {
		relayErr = common.WrapperErr(err, common.ErrCodeInvalidRequest, 400)
		return
	}
	common.Logger.Debugw("generateSong get", "requestBody", string(reqBody))
	// check params
	if params.Prompt == "" {
		relayErr = common.WrapperErr(fmt.Errorf("prompt is empty"), common.ErrCodeInvalidRequest, 400)
		return
	}

	requestBody, _ := json.Marshal(params)

	resp, err := doRequest("POST", common.BaseUrl+"/api/generate/lyrics/", bytes.NewReader(requestBody), map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + AccountInst.Certificate.JWT,
	})
	if err != nil {
		relayErr = common.WrapperErr(err, common.ErrCodeInternalError, 500)
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
		common.Logger.Errorw("generateSong", "status", resp.StatusCode, "responseBody", string(responseBody))
		relayErr = common.WrapperErr(fmt.Errorf(resp.Status), common.ErrCodeInvalidRequest, resp.StatusCode)
		return
	}
	common.Logger.Debugw("generateSong", "responseBody", string(responseBody))

	var sunoResp GenLyricsResponse
	err = json.Unmarshal(responseBody, &sunoResp)
	if err != nil {
		relayErr = common.WrapperErr(err, common.ErrCodeInternalError, 500)
		return
	}

	var lyricsResp FetchLyricsResponse

	relayError := common.LoopTask(func() (bool, *common.RelayError) {
		resp, err = doRequest("GET", fmt.Sprintf(common.BaseUrl+"/api/generate/lyrics/%s", sunoResp.ID), nil, map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer " + AccountInst.Certificate.JWT,
		})
		if err != nil {
			common.Logger.Errorw("请求获取歌词失败", "err", err)
			return false, common.WrapperErr(err, common.ErrCodeInternalError, 500)
		}
		responseBody, _ = io.ReadAll(resp.Body)
		if resp.StatusCode != 200 {
			common.Logger.Errorw("请求获取歌词失败", "resp.StatusCode", resp.StatusCode, "responseBody", string(responseBody))
			return false, nil
		}

		defer func(Body io.ReadCloser) {
			err = Body.Close()
			if err != nil {
				common.Logger.Errorw("body close", "err", err)
			}
		}(resp.Body)

		err = json.Unmarshal(responseBody, &lyricsResp)
		if err != nil {
			return false, nil
		}
		common.Logger.Debugw("生成歌词", "responseBody", string(responseBody))
		if lyricsResp.Status == "complete" || lyricsResp.Text != "" {
			return true, nil
		}
		return false, nil

	}, 5*time.Second, 2*time.Minute, false)
	if relayError != nil {
		relayErr = relayError
		return
	}

	task := po.TaskWithData[po.SunoLyrics]{
		Action:     po.TaskActionLyrics,
		TaskID:     sunoResp.ID,
		Status:     readerStatus(lyricsResp.Status),
		SubmitTime: time.Now().Unix(),
		StartTime:  time.Now().Unix(),
		FinishTime: time.Now().Unix(),
		Data: po.SunoLyrics{
			ID:     sunoResp.ID,
			Status: lyricsResp.Status,
			Title:  lyricsResp.Title,
			Text:   lyricsResp.Text,
		},
	}
	// 保存任务
	err = po.AddTasks([]po.TaskWithData[po.SunoLyrics]{task})
	if err != nil {
		relayErr = common.WrapperErr(err, common.ErrCodeInternalError, 500)
		return
	}
	// writer response
	//c.JSON(200, ginplus.ApiRetSucc(task.TaskID))
	return task.TaskID, nil
}

func loopFetchTask(taskID string) (bool, *common.RelayError) {

	oldTask, exist, err := po.GetTaskByTaskID[po.SunoSongs](taskID)
	if err != nil {
		return true, common.WrapperErr(err, common.ErrCodeInternalError, 500)
	}
	if !exist {
		return true, common.WrapperErr(fmt.Errorf("task %s not found", taskID), common.ErrCodeInternalError, 500)
	}
	if oldTask.Status.IsDone() {
		return true, nil
	}
	timeoutAt := time.Now().Add(-time.Duration(common.ChatTimeOut) * time.Second).Unix()
	if oldTask.SubmitTime > 0 && timeoutAt > oldTask.SubmitTime {
		err = oldTask.TaskFailed("time out")
		return true, common.WrapperErr(fmt.Errorf("task %s time out", taskID), common.ErrCodeInternalError, 500)
	}

	ids := lo.Map(oldTask.Data, func(data po.SunoSong, idx int) string {
		return data.ID
	})

	common.Logger.Debugw("loopFetchTask", "ids", ids)
	if len(ids) == 0 {
		return true, nil
	}
	resp, err := doRequest("GET", fmt.Sprintf(common.BaseUrl+"/api/feed/?ids=%s", strings.Join(ids, ",")), nil, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + AccountInst.Certificate.JWT,
	})
	if err != nil {
		return false, common.WrapperErr(err, common.ErrCodeInternalError, 500)
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			common.Logger.Errorw("body close", "err", err)
		}
	}(resp.Body)
	responseBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		common.Logger.Errorw("轮询返回失败", "status code", resp.Status, "responseBody", string(responseBody))
		return false, nil
	}
	common.Logger.Debugw("loopFetchTask", "responseBody", string(responseBody))

	var clips []po.SunoSong
	err = json.Unmarshal(responseBody, &clips)
	task := po.TaskWithData[po.SunoSongs]{
		BaseModel:  oldTask.BaseModel,
		TaskID:     oldTask.TaskID,
		Action:     oldTask.Action,
		Status:     oldTask.Status,
		FailReason: oldTask.FailReason,
		SubmitTime: oldTask.SubmitTime,
		StartTime:  oldTask.StartTime,
		FinishTime: oldTask.FinishTime,
	}
	task.Data = clips

	var status string
	if len(clips) == 2 {
		if clips[0].Status == clips[1].Status {
			status = clips[0].Status
		} else {
			if clips[0].Status == "complete" {
				status = clips[1].Status
			} else {
				status = clips[0].Status
			}
		}
	} else {
		status = clips[0].Status
	}

	task.Status = readerStatus(status)

	if task.Status == po.TaskStatusInProgress && task.StartTime == 0 {
		task.StartTime = time.Now().Unix()
	}
	if task.FinishTime == 0 && task.Status.IsDone() {
		task.FinishTime = time.Now().Unix()
		if task.Status == po.TaskStatusFailure {
			for _, c := range clips {
				msg := common.Any2String(c.Metadata.ErrorMessage)
				if msg != "" {
					task.FailReason = msg
					break
				}
			}
		}
	}

	if po.CheckTaskNeedUpdate(oldTask, &task) {
		// update task
		common.Logger.Debugw("loopFetchTask UpdateTask", "task", task)
		err = po.UpdateTask(task.ID, &task)
		if err != nil {
			return false, common.WrapperErr(err, common.ErrCodeInternalError, 500)
		}
	}

	if task.FinishTime != 0 {
		return true, nil
	}
	return false, nil
}

// @Summary Fetch task
// @Schemes
// @Description
// @Tags suno
// @Accept json
// @Produce json
// @Param body body FetchReq true "fetch task ids"
// @Success 200 {object} ginplus.DataResult{data=[]po.Task} "song tasks"
// @Router /api/fetch [post]
func Fetch(c *gin.Context) {
	var params FetchReq
	err := common.UnmarshalBodyReusable(c, &params)
	if err != nil {
		common.ReturnErr(c, err, common.ErrCodeInvalidRequest, 400)
		return
	}

	tasks, err := po.GetTaskByTaskIDs[json.RawMessage](params.IDs)
	if err != nil {
		common.ReturnErr(c, err, common.ErrCodeInternalError, 500)
		return
	}

	c.JSON(200, ginplus.ApiRetSucc(tasks))
}

// @Summary Fetch task by id
// @Schemes
// @Description
// @Tags suno
// @Accept json
// @Produce json
// @Param id path string true "fetch single task by id"
// @Success 200 {object} ginplus.DataResult{data=po.Task} "song task"
// @Router /api/fetch/{id} [get]
func FetchByID(c *gin.Context) {

	id := c.Param("id")
	if id == "" {
		common.ReturnErr(c, fmt.Errorf("invalid id"), common.ErrCodeInvalidRequest, 400)
		return
	}

	tasks, exist, err := po.GetTaskByTaskID[json.RawMessage](id)
	if err != nil {
		common.ReturnErr(c, err, common.ErrCodeInternalError, 500)
		return
	}

	if !exist {
		common.ReturnErr(c, fmt.Errorf("task not exist"), common.ErrCodeInvalidRequest, 400)
		return
	}

	c.JSON(200, ginplus.ApiRetSucc(tasks.Data))
}

// @Summary Get Account config
// @Schemes
// @Description
// @Tags account
// @Accept json
// @Produce json
// @Success 200 {object} ginplus.DataResult{data=Account} "song task"
// @Router /api/account [get]
func GetAccount(c *gin.Context) {
	err := getCredits(AccountInst)
	if err != nil {
		common.ReturnErr(c, err, common.ErrCodeInternalError, 500)
		return
	}

	c.JSON(200, ginplus.ApiRetSucc(AccountInst))
}

func readerStatus(status string) po.TaskStatus {
	switch status {
	case "submitted":
		return po.TaskStatusSubmitted
	case "queued":
		return po.TaskStatusQueued
	case "streaming":
		return po.TaskStatusInProgress
	case "complete":
		return po.TaskStatusSuccess
	case "error":
		return po.TaskStatusFailure
	default:
		common.Logger.Errorw("unknow status: " + status)
		return po.TaskStatusFailure
	}
}
