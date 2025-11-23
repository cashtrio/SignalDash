package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/developerasun/SignalDash/server/config"
	"github.com/developerasun/SignalDash/server/controller"
	docs "github.com/developerasun/SignalDash/server/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title SignalDash API
// @version 1.0
// @description SignalDash backend API documentation
// @BasePath /
func main() {
	options := &config.ViperOptions{
		Filename:  "options",
		ConfigDir: "config",
	}
	v := options.InitConfig()
	port := v.GetString("server.port")

	apiServer := gin.Default()
	apiServer.SetTrustedProxies(nil)

	docs.SwaggerInfo.BasePath = ""
	apiServer.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// TODO move to middleware
	apiServer.Use(ErrorHandler())
	apiServer.Use(gin.Recovery())

	apiServer.GET("/api/health", controller.Health)

	apiServer.Run(":" + port)
	log.Println("main.go: router started")
}

// ErrorHandler captures errors and returns a consistent JSON error response
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // Step1: Process the request first.

		// Step2: Check if any errors were added to the context
		if len(c.Errors) > 0 {
			// Step3: Use the last error
			err := c.Errors.Last().Err

			fmt.Println(err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "error!",
			})
		}

		// Any other steps if no errors are found
	}
}
