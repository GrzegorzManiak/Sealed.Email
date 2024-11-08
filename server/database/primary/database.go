package primary

import (
	models2 "github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitiateConnection() *gorm.DB {
	databaseConnection, err := gorm.Open(sqlite.Open("./dev/primary.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
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
