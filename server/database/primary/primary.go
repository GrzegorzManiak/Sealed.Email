package primary

import (
	models "github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
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
	databaseConnection, err := gorm.Open(sqlite.Open(dbPath+"primary.db"), &gorm.Config{Logger: logger})
	if err != nil {
		zap.L().Panic("failed to connect to database", zap.Error(err))
	}

	AutoMigrateTables(databaseConnection)
	return databaseConnection
}

func AutoMigrateTables(databaseConnection *gorm.DB) {
	databaseConnection.AutoMigrate(&models.UserVerify{})
	databaseConnection.AutoMigrate(&models.Session{})
	databaseConnection.AutoMigrate(&models.User{})
	databaseConnection.AutoMigrate(&models.UserDomain{})
	databaseConnection.AutoMigrate(&models.UserHold{})
	databaseConnection.AutoMigrate(&models.UserQuota{})
	databaseConnection.AutoMigrate(&models.UserEmail{})
}
