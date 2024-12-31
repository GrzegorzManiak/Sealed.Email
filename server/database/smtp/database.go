package smtp

import (
	"github.com/GrzegorzManiak/NoiseBackend/database/smtp/models"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
	"os"
)

func InitiateConnection() *gorm.DB {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./dev/"
	}

	logger := zapgorm2.New(zap.L())
	logger.SetAsDefault()
	driver := sqlite.Open(dbPath + "smtp.db")
	databaseConnection, err := gorm.Open(driver, &gorm.Config{Logger: logger})
	if err != nil {
		zap.L().Panic("failed to connect to database", zap.Error(err))
	}

	AutoMigrateTables(databaseConnection)
	return databaseConnection
}

func AutoMigrateTables(databaseConnection *gorm.DB) {
	databaseConnection.AutoMigrate(&models.InboundEmail{})
}
