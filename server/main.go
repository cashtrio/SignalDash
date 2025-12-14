package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/cashtrio/SignalDash/server/config"
	"github.com/cashtrio/SignalDash/server/instance"
)

// @title SignalDash API
// @version 1.0
// @description SignalDash backend API documentation
// @BasePath /
func main() {
	environment := config.NewEnvironment("config", "options")
	port := environment.Instance.GetString("server.port")
	databaseName := environment.Instance.GetString("server.database.main")

	database := instance.NewDatabase(databaseName).DB
	apiServer := instance.NewApiServer(gin.Default(), database)
	cronWorker := instance.NewCronWorker(environment.Instance)
	cronWorker.Run()
	apiServer.Run(":" + port)

	log.Println("main.go: router started")
}
