package worker

import (
	"github.com/GrzegorzManiak/NoiseBackend/database/smtp/models"
	"gorm.io/gorm"
)

func getEmailById(emailId string, queueDatabaseConnection *gorm.DB) (*models.OutboundEmail, error) {
	email := &models.OutboundEmail{}
	err := queueDatabaseConnection.
		Preload("OutboundEmailKeys").
		Where("email_id = ?", emailId).
		First(email).Error

	if err != nil {
		return nil, err
	}
	return email, nil
}
