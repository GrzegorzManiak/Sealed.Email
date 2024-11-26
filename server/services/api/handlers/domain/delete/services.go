package domainDelete

import (
	"github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"gorm.io/gorm"
)

func deleteDomain(
	userID uint,
	domainID string,
	databaseConnection *gorm.DB,
) helpers.AppError {
	domain := &models.UserDomain{}
	dbQuery := databaseConnection.Where("user_id = ? AND r_id = ?", userID, domainID).First(domain)
	if dbQuery.Error != nil {
		return helpers.GenericError{
			Message: dbQuery.Error.Error(),
			ErrCode: 400,
		}
	}

	dbDelete := databaseConnection.Delete(domain)
	if dbDelete.Error != nil {
		return helpers.GenericError{
			Message: dbDelete.Error.Error(),
			ErrCode: 400,
		}
	}

	return nil
}
