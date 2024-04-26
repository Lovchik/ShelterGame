package sqlite

import (
	"ShelterGame/internal/config"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var database *gorm.DB

func init() {
	db, err := gorm.Open(sqlite.Open(config.GetConfig().DatabaseUrl), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database")
	}
	database = db
}

func GetDB() *gorm.DB {
	return database
}
