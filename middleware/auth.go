package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"suno-api/common"
)

func SecretAuth() func(c *gin.Context) {
	return func(c *gin.Context) {
		if common.SecretToken == "" {
			return
		}
		accessToken := c.Request.Header.Get("Authorization")
		accessToken = strings.TrimLeft(accessToken, "Bearer ")
		if accessToken == common.SecretToken {
			c.Next()
		} else {
			common.ReturnErr(c, fmt.Errorf("unauthorized secret token"), common.ErrCodeInvalidRequest, http.StatusUnauthorized)
			c.Abort()
			return
		}
	}
}
