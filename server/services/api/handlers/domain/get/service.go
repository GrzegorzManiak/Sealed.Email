package domainGet

import (
	"github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/errors"
	"gorm.io/gorm"
)

func getDomain(
	user *models.User,
	domainID string,
	databaseConnection *gorm.DB,
) (*models.UserDomain, errors.AppError) {
	var domain models.UserDomain
	result := databaseConnection.
		Where("p_id = ? AND user_id = ?", domainID, user.ID).
		First(&domain)

	if result.Error != nil {
		return &models.UserDomain{}, errors.NotFound(
			"We could not find the domain you are looking for. Please try again.",
			"Domain not found!",
		)
	}

	return &domain, nil
}
