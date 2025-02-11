package server

import (
	"github.com/GrzegorzManiak/NoiseBackend/database/smtp/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/queue"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func getEmailById(emailId string, queueDatabaseConnection *gorm.DB) (*models.InboundEmail, error) {
	email := &models.InboundEmail{}
	err := queueDatabaseConnection.
		Where("ref_id = ?", emailId).
		First(email).Error

	if err != nil {
		return nil, err
	}
	return email, nil
}

func fetchRecipients(email *models.InboundEmail) []string {
	toArray := make([]string, 0, len(email.To))
	for k := range email.To {
		toArray = append(toArray, k)
	}
	return toArray
}

func Worker(entry *queue.Entry, queueDatabaseConnection *gorm.DB) int8 {
	email, err := getEmailById(entry.RefID, queueDatabaseConnection)
	if err != nil {
		zap.L().Debug("Failed to get email by id", zap.Error(err))
		return 1
	}

	zap.L().Debug("Processing email", zap.String("email_id", email.EmailId))
	return 2
}
