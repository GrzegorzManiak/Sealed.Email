package get

import (
	"github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/errors"
	"gorm.io/gorm"
)

func fetchEmail(
	user *models.User,
	input Input,
	databaseConnection *gorm.DB,
) (*models.UserEmail, errors.AppError) {
	email := &models.UserEmail{}
	if err := databaseConnection.
		Table("user_emails").
		Select([]string{"read", "folder", "p_id", "domain_p_id", "`to`", "received_at", "sent", "bucket_path"}).
		Where("user_id = ? AND domain_p_id = ? AND bucket_path = ?", user.ID, input.DomainID, input.BucketPath).
		First(email).Error; err != nil {
		return nil, errors.User("The requested email could not be found.", "Email not found!")
	}

	return email, nil
}
