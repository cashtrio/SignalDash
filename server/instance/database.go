package instance

import (
	"log"

	"github.com/cashtrio/SignalDash/server/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type database struct {
	DB *gorm.DB
}

func NewDatabase(databaseName string) *database {
	db, oErr := gorm.Open(sqlite.Open(databaseName), &gorm.Config{})
	if oErr != nil {
		log.Fatalf("main.go: failed to open sqlite")
	}
	db.AutoMigrate(&models.Indicator{})
	log.Println("main.go: database migrated and opened")

	return &database{
		DB: db,
	}
}
