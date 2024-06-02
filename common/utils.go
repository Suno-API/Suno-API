package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

func SafeGoroutine(f func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				Logger.Error(fmt.Sprintf("child goroutine panic occured: error: %v, stack: %s", r, string(debug.Stack())))
			}
		}()
		f()
	}()
}

func Any2String(inter interface{}) string {
	if inter == nil {
		return ""
	}
	switch inter.(type) {
	case string:
		return inter.(string)
	case *string:
		ptr, ok := inter.(*string)
		if !ok || ptr == nil {
			return ""
		}
		return *ptr
	case int, int64:
		return fmt.Sprintf("%d", inter)
	case float64:
		return fmt.Sprintf("%f", inter)
	}
	return "Not Implemented"
}

func Any2Int(data any) int {
	if data == nil {
		return 0
	}
	if valueAssert, ok := data.(float64); ok {
		return int(valueAssert)
	}
	if valueAssert, ok := data.(int64); ok {
		return int(valueAssert)
	}
	if valueAssert, ok := data.(int); ok {
		return valueAssert
	}
	if valueAssert, ok := data.(*int); ok {
		return *valueAssert
	}
	if valueAssert, ok := data.(string); ok {
		valueAssertInt, _ := strconv.Atoi(valueAssert)
		return valueAssertInt
	}
	return 0
}

func Any2Bool(data any) bool {
	if data == nil {
		return false
	}
	switch data.(type) {
	case string:
		if data == "true" {
			return true
		}
		if data == "1" {
			return true
		}
		return false
	case int:
		if data == 1 {
			return true
		}
		return false
	case bool:
		return data.(bool)
	case nil:
		return false
	}
	return false
}

func GetUUID() string {
	code := uuid.New().String()
	code = strings.Replace(code, "-", "", -1)
	return code
}

const keyChars = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GetRandomString(length int) string {
	//rand.Seed(time.Now().UnixNano())
	key := make([]byte, length)
	for i := 0; i < length; i++ {
		key[i] = keyChars[rand.Intn(len(keyChars))]
	}
	return string(key)
}

func GetTimeString() string {
	now := time.Now()
	return fmt.Sprintf("%s%d", now.Format("20060102150405"), now.UnixNano()%1e9)
}

func GetOrDefault(env string, defaultValue int) int {
	if env == "" || os.Getenv(env) == "" {
		return defaultValue
	}
	num, err := strconv.Atoi(os.Getenv(env))
	if err != nil {
		Logger.Error(fmt.Sprintf("failed to parse %s: %s, using default value: %d", env, err.Error(), defaultValue))
		return defaultValue
	}
	return num
}

func GetOrDefaultString(env string, defaultValue string) string {
	if env == "" || os.Getenv(env) == "" {
		return defaultValue
	}
	return os.Getenv(env)
}

func GetJsonString(data any) string {
	if data == nil {
		return ""
	}
	b, _ := json.Marshal(data)
	return string(b)
}

func UnmarshalBodyReusable(c *gin.Context, v any) error {
	requestBody, err := GetRequestBody(c)
	if err != nil {
		return err
	}
	contentType := c.Request.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		err = json.Unmarshal(requestBody, &v)
	} else {
		// skip for now
		// try
		err = json.Unmarshal(requestBody, &v)
	}
	if err != nil {
		return err
	}
	// Reset request body
	c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
	return nil
}

const KeyRequestBody = "key_request_body"

func GetRequestBody(c *gin.Context) ([]byte, error) {
	requestBody, _ := c.Get(KeyRequestBody)
	if requestBody != nil {
		return requestBody.([]byte), nil
	}
	requestBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil, err
	}
	_ = c.Request.Body.Close()
	c.Set(KeyRequestBody, requestBody)
	return requestBody.([]byte), nil
}

func LoopTask(task func() (bool, *RelayError), duration time.Duration, timeout time.Duration, sync bool) *RelayError {

	fn := func() *RelayError {
		for {
			select {
			case <-time.After(duration):
				stop, err := task()
				if err != nil {
					return err
				}
				if stop {
					return nil
				}
			case <-time.After(timeout):
				Logger.Errorw("Operation timed out after", "timeout", timeout)
				return nil
			}
		}
	}
	if sync {
		SafeGoroutine(func() {
			err := fn()
			if err != nil {
				Logger.Errorw("loop task failed", "err", err)
			}
		})
	} else {
		err := fn()
		return err
	}
	return nil
}

func GetRootDir() string {
	_, filename, _, _ := runtime.Caller(0)
	utilDir := filepath.Dir(filename)
	paths, _ := filepath.Split(utilDir)
	return paths
}

func ToPtr[T any](arg T) *T {
	return &arg
}
