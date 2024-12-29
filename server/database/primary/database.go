package primary

import (
	models2 "github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
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
	databaseConnection.AutoMigrate(&models2.UserVerify{})
	databaseConnection.AutoMigrate(&models2.Session{})
	databaseConnection.AutoMigrate(&models2.User{})
	databaseConnection.AutoMigrate(&models2.UserDomain{})
	databaseConnection.AutoMigrate(&models2.UserInbox{})
	databaseConnection.AutoMigrate(&models2.UserHold{})
	databaseConnection.AutoMigrate(&models2.UserQuota{})
}
