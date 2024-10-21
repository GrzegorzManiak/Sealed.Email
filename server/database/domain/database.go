package domain

import (
	"github.com/GrzegorzManiak/NoiseBackend/database/domain/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitiateConnection() (*gorm.DB, gorm.Dialector) {
	driver := sqlite.Open("./dev/domain.db")
	databaseConnection, err := gorm.Open(driver, &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	AutoMigrateTables(databaseConnection)
	return databaseConnection, driver
}

func AutoMigrateTables(databaseConnection *gorm.DB) {
	databaseConnection.AutoMigrate(&models.QueueEntry{})
}
