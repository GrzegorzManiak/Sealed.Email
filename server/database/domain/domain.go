package domain

import (
	"os"

	"github.com/GrzegorzManiak/NoiseBackend/database/domain/models"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
)

func InitiateConnection() *gorm.DB {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./dev/"
	}

	logger := zapgorm2.New(zap.L())
	logger.SetAsDefault()

	driver := sqlite.Open(dbPath + "domain.db")

	databaseConnection, err := gorm.Open(driver, &gorm.Config{Logger: logger})
	if err != nil {
		zap.L().Panic("failed to connect to database", zap.Error(err))
	}

	AutoMigrateTables(databaseConnection)

	return databaseConnection
}

func AutoMigrateTables(databaseConnection *gorm.DB) {
	databaseConnection.AutoMigrate(&models.VerificationQueue{})
}
