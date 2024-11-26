package domainVerify

import (
	"github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"gorm.io/gorm"
)

func fetchDomainByID(
	userID string,
	domainID string,
	databaseConnection *gorm.DB,
) (*models.UserDomain, helpers.AppError) {
	domain := &models.UserDomain{}
	dbQuery := databaseConnection.Where("user_id = ? AND r_id = ?", userID, domainID).First(domain)
	if dbQuery.Error != nil {
		return nil, helpers.GenericError{
			Message: dbQuery.Error.Error(),
			ErrCode: 400,
		}
	}

	return domain, nil
}
