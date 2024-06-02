package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"runtime/debug"
	"suno-api/common"
	"suno-api/entity/po"
	"suno-api/middleware"
)

func main() {
	common.SetupLogger()
	common.Logger.Info("Suno-API started")

	if common.DebugEnabled {
		common.Logger.Info("running in debug mode")
		gin.SetMode(gin.ReleaseMode)
	}
	if common.PProfEnabled {
		common.SafeGoroutine(func() {
			log.Println(http.ListenAndServe("0.0.0.0:8005", nil))
		})
		common.Logger.Info("running in pprof")
	}

	err := po.InitDB()
	if err != nil {
		common.Logger.Fatal("failed to initialize database: " + err.Error())
	}
	// Initialize HTTP server
	server := gin.New()
	server.Use(middleware.RequestId())
	server.Use(gin.CustomRecovery(func(c *gin.Context, err any) {
		common.Logger.Error(fmt.Sprintf("panic detected: %v, stack: %s", err, string(debug.Stack())))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"message": fmt.Sprintf("Panic detected, error: %v. Please contact site admin", err),
				"type":    "api_panic",
			},
		})
	}))
	server.Use(middleware.GinzapWithConfig())

	InitService()
	RegisterRouter(server)

	common.Logger.Info("Start: 0.0.0.0:" + common.Port)
	err = server.Run(":" + common.Port)
	if err != nil {
		common.Logger.Fatal("failed to start HTTP server: " + err.Error())
	}
}
