package domainDelete

import (
	"github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func deleteDomain(
	userID uint,
	domainID string,
	databaseConnection *gorm.DB,
) errors.AppError {
	domain := &models.UserDomain{}
	dbQuery := databaseConnection.Where("user_id = ? AND p_id = ?", userID, domainID).First(domain)
	if dbQuery.Error != nil {
		zap.L().Debug("Error querying domain", zap.Error(dbQuery.Error), zap.Any("domain", domain))
		return errors.User("The requested domain could not be found.", "Domain not found!")
	}

	dbDelete := databaseConnection.Delete(domain)
	if dbDelete.Error != nil {
		zap.L().Debug("Error deleting domain", zap.Error(dbDelete.Error), zap.Any("domain", domain))
		return errors.User("The requested domain could not be deleted.", "Failed to delete domain!")
	}

	return nil
}
