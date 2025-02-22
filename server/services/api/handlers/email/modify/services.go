package modify

import (
	"github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/errors"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/services"
	"gorm.io/gorm"
	"strings"
)

func createUpdateQuery(
	input Input,
	query *gorm.DB,
) {
	// -- Read
	input.Read = strings.ToLower(input.Read)
	if input.Read == "read" {
		query = query.Update("read", 1)
	} else if input.Read == "unread" {
		query = query.Update("read", 0)
	}

	// -- Folder
	if input.Folder != "" {
		query = query.Update("folder", input.Folder)
	}

	// -- Spam
	input.Spam = strings.ToLower(input.Spam)
	if input.Spam == "true" {
		query = query.Update("spam", 1)
	} else if input.Spam == "false" {
		query = query.Update("spam", 0)
	}

	// -- Pinned
	input.Pinned = strings.ToLower(input.Pinned)
	if input.Pinned == "true" {
		query = query.Update("pinned", 1)
	} else if input.Pinned == "false" {
		query = query.Update("pinned", 0)
	}
}

func updateEmail(
	input *Input,
	data *services.Handler,
) errors.AppError {

	query := data.DatabaseConnection.Model(&models.UserEmail{}).
		Where("user_id = ? AND p_id in (?)", data.User.ID, input.EmailIds)

	createUpdateQuery(*input, query)

	if err := query.Error; err != nil {
		return errors.Server(
			"Could not update the email.",
			"Failed to update email",
		)
	}

	return nil
}
