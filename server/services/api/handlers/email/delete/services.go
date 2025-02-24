package delete

import (
	"github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/errors"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/services"
)

func deleteEmails(
	input *Input,
	data *services.Handler,
) errors.AppError {
	query := data.DatabaseConnection.
		Where("user_id = ? AND p_id in (?)", data.User.ID, input.EmailIds).
		Delete(&models.UserEmail{})

	if err := query.Error; err != nil {
		return errors.Server(
			"Sorry! For some reason we could not delete the emails",
			"Failed to delete emails",
		)
	}

	return nil
}
