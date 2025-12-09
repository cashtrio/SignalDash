package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/developerasun/SignalDash/server/config"
	"github.com/developerasun/SignalDash/server/instance"
	"github.com/developerasun/SignalDash/server/models"
)

// @title SignalDash API
// @version 1.0
// @description SignalDash backend API documentation
// @BasePath /
func main() {
	environment := config.NewEnvironment("config", "options")
	port := environment.Instance.GetString("server.port")
	databaseName := environment.Instance.GetString("server.database.main")

	db, oErr := gorm.Open(sqlite.Open(databaseName), &gorm.Config{})
	if oErr != nil {
		log.Fatalf("main.go: failed to open sqlite")
	}
	db.AutoMigrate(&models.Indicator{})
	log.Println("main.go: database migrated and opened")

	apiServer := instance.NewApiServer(gin.Default(), db)
	cronWorker := instance.NewCronWorker(environment.Instance)
	cronWorker.Run()
	apiServer.Run(":" + port)

	log.Println("main.go: router started")
}
