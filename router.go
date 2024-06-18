package main

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	docs "suno-api/docs"
)

func RegisterRouter(r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	apiRouter := r.Group("/suno")
	{
		apiRouter.POST("/submit/:action", Submit)
		apiRouter.GET("/fetch/:id", FetchByID)
		apiRouter.POST("/fetch", Fetch)

		apiRouter.GET("/account", GetAccount)
	}
	// chat
	r.POST("/v1/chat/completions", ChatCompletions)
}
