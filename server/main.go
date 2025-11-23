package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/developerasun/SignalDash/server/config"
	"github.com/developerasun/SignalDash/server/controller"
	docs "github.com/developerasun/SignalDash/server/docs"
	"github.com/developerasun/SignalDash/server/models"
	"github.com/developerasun/SignalDash/server/repository"
	"github.com/developerasun/SignalDash/server/service"
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

	db, oErr := gorm.Open(sqlite.Open("sd_app.db"), &gorm.Config{})
	if oErr != nil {
		log.Fatalf("main.go: failed to open sqlite")
	}
	db.AutoMigrate(&models.Indicator{})
	log.Println("main.go: database opened")

	indicatorRepo := repository.NewIndicator(db)
	indicatorService := service.NewIndicatorService(indicatorRepo)
	indicatorController := controller.NewIndicatorController(indicatorService)

	docs.SwaggerInfo.BasePath = ""
	apiServer.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// TODO move to middleware
	apiServer.Use(ErrorHandler())
	apiServer.Use(gin.Recovery())

	apiServer.GET("/api/health", controller.Health)
	apiServer.GET("/api/indicator", indicatorController.GetIndicator)

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
