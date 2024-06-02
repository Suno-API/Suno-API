package common

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"suno-api/lib/ginplus"
)

type RelayError struct {
	StatusCode int
	Code       string
	Err        error
	LocalErr   bool
}

const (
	ErrCodeInvalidRequest = "invalid_request"
	ErrCodeInternalError  = "internal_error"
)

func ReturnErr(c *gin.Context, err error, code string, statusCode int) {
	c.JSON(statusCode, ginplus.BuildApiReturn(code, err.Error(), nil))
}

func ReturnRelayErr(c *gin.Context, relayErr *RelayError) {
	if relayErr.Err == nil {
		relayErr.Err = fmt.Errorf("unknown error")
	}
	c.JSON(relayErr.StatusCode, ginplus.BuildApiReturn(relayErr.Code, relayErr.Err.Error(), nil))
}

func WrapperErr(err error, code string, statusCode int) *RelayError {
	return &RelayError{
		StatusCode: statusCode,
		Code:       code,
		Err:        err,
	}
}
