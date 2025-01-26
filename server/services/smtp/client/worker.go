package client

import (
	"crypto/tls"
	"github.com/GrzegorzManiak/NoiseBackend/database/smtp/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/queue"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func getEmailById(emailId string, queueDatabaseConnection *gorm.DB) (*models.OutboundEmail, error) {
	email := &models.OutboundEmail{}
	err := queueDatabaseConnection.Where("email_id = ?", emailId).First(email).Error
	if err != nil {
		return nil, err
	}
	return email, nil
}

func Worker(certs *tls.Config, entry *queue.Entry, queueDatabaseConnection *gorm.DB) int8 {
	zap.L().Debug("Processing smtp queue", zap.Any("entry", entry))

	emailId, err := models.UnmarshalQueueEmailId(entry.Data)
	if err != nil {
		zap.L().Debug("Failed to unmarshal email id", zap.Error(err))
		return 2
	}

	email, err := getEmailById(emailId.EmailId, queueDatabaseConnection)
	if err != nil {
		zap.L().Debug("Failed to get email by id", zap.Error(err))
		return 2
	}

	for _, to := range email.To {
		if err := attemptSendEmail(certs, email, to); err != nil {
			zap.L().Debug("Failed to send email", zap.Error(err))
			return 2
		}
	}

	return 1
}
