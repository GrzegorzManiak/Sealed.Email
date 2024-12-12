package domainVerify

import (
	"github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"gorm.io/gorm"
)

func fetchDomainByID(
	userID uint,
	domainID string,
	databaseConnection *gorm.DB,
) (*models.UserDomain, helpers.AppError) {
	domain := &models.UserDomain{}
	dbQuery := databaseConnection.Where("user_id = ? AND p_id = ?", userID, domainID).First(domain)
	if dbQuery.Error != nil {
		return nil, helpers.NewServerError("The requested domain could not be found.", "Domain not found!")
	}

	return domain, nil
}
