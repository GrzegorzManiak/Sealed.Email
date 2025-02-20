package services

import (
	"github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/errors"
	"gorm.io/gorm"
	"sync"
)

var primaryDatabaseMutex = sync.Mutex{}

func VerifyDomain(domain string, uid uint64, did uint64, databaseConnection *gorm.DB) errors.AppError {
	primaryDatabaseMutex.Lock()
	defer primaryDatabaseMutex.Unlock()

	err := databaseConnection.Transaction(func(tx *gorm.DB) error {

		// -- Mark all domains as unverified
		tx.Model(&models.UserDomain{}).Where("domain = ? AND id != ?", domain, did).Update("verified", false)

		// -- Mark the domain as verified
		tx.Model(&models.UserDomain{}).Where("domain = ? AND user_id = ? AND id = ?", domain, uid, did).Update("verified", true)

		return nil
	})

	if err != nil {
		return errors.NotFound(
			"We were unable to verify the domain. Please try again.",
			"Failed to verify domain!",
		)
	}

	return nil
}
