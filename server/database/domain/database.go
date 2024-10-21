package domain

import (
	"github.com/GrzegorzManiak/NoiseBackend/database/domain/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitiateConnection() *gorm.DB {
	databaseConnection, err := gorm.Open(sqlite.Open("./dev/domain.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	AutoMigrateTables(databaseConnection)
	return databaseConnection
}

func AutoMigrateTables(databaseConnection *gorm.DB) {
	databaseConnection.AutoMigrate(&models.QueueEntry{})
}
