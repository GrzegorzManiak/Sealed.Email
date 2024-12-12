package add

import (
	"github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"gorm.io/gorm"
)

func getDomain(
	user *models.User,
	domainID string,
	databaseConnection *gorm.DB,
) (*models.UserDomain, helpers.AppError) {
	var domain models.UserDomain
	result := databaseConnection.
		Where("p_id = ? AND user_id = ?", domainID, user.ID).
		First(&domain)

	if result.Error != nil {
		return &models.UserDomain{}, helpers.NewNotFoundError(
			"We could not find the domain you are looking for. Please try again.",
			"Domain not found!",
		)
	}

	return &domain, nil
}

func createInbox(
	user *models.User,
	domain *models.UserDomain,
	input *Input,
	databaseConnection *gorm.DB,
) (*models.UserInbox, helpers.AppError) {
	inbox := models.UserInbox{
		UserID:    user.ID,
		DomainID:  domain.ID,
		EmailHash: input.EmailHash,

		AsymmetricPrivateKey: input.AsymmetricPrivateKey,
		AsymmetricPublicKey:  input.AsymmetricPublicKey,
		SymmetricRootKey:     input.SymmetricRootKey,

		Version: 1,
	}

	if err := databaseConnection.Create(&inbox); err.Error != nil {
		return &models.UserInbox{}, helpers.NewServerError(
			"Inbox could not be created. Please contact support if this issue persists.",
			"Failed to create inbox!",
		)
	}

	return &inbox, nil
}
