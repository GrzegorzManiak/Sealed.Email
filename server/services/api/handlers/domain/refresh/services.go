package domainVerify

import (
	"github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/errors"
	"gorm.io/gorm"
)

func fetchDomainByID(
	userID uint,
	domainID string,
	databaseConnection *gorm.DB,
) (*models.UserDomain, errors.AppError) {
	domain := &models.UserDomain{}

	dbQuery := databaseConnection.Where("user_id = ? AND p_id = ?", userID, domainID).First(domain)
	if dbQuery.Error != nil {
		return nil, errors.Server("The requested domain could not be found.", "Domain not found!")
	}

	return domain, nil
}
