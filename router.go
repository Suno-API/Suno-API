package main

import (
	"github.com/gin-gonic/gin"
)

type SunoAPI struct {
}

func (s *SunoAPI) Name() string {
	return "suno"
}

var Service = &SunoAPI{}

func RegisterRouter(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	apiRouter := r.Group("/api")
	{
		apiRouter.POST("/submit/:action", Submit)
		apiRouter.GET("/fetch/:id", FetchByID)
		apiRouter.POST("/fetch", Fetch)

		apiRouter.GET("/account", GetAccount)

		// chat
		apiRouter.POST("/v1/chat/completions", ChatCompletions)
	}
}
