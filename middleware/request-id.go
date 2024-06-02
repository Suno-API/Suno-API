package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"suno-api/common"
)

const RequestIdKey = "X-Any2api-Request-Id"

func RequestId() func(c *gin.Context) {
	return func(c *gin.Context) {
		id := common.GetTimeString() + common.GetRandomString(8)
		c.Set(RequestIdKey, id)
		ctx := context.WithValue(c.Request.Context(), RequestIdKey, id)
		c.Request = c.Request.WithContext(ctx)
		c.Header(RequestIdKey, id)
		c.Next()
	}
}
