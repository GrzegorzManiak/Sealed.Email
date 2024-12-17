package inboxGet

import (
	"github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"gorm.io/gorm"
)

func getInbox(
	user *models.User,
	inboxID string,
	databaseConnection *gorm.DB,
) (*models.UserInbox, helpers.AppError) {
	var inbox models.UserInbox
	result := databaseConnection.
		Where("p_id = ? AND user_id = ?", inboxID, user.ID).
		First(&inbox)

	if result.Error != nil {
		return &models.UserInbox{}, helpers.NewNotFoundError(
			"We could not find the inbox that you are looking for. Please try again.",
			"Inbox not found!",
		)
	}

	return &inbox, nil
}
