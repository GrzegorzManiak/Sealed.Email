package server

import (
	"crypto/tls"
	"github.com/GrzegorzManiak/NoiseBackend/database/smtp/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/queue"
	"gorm.io/gorm"
)

func getEmailById(emailId string, queueDatabaseConnection *gorm.DB) (*models.InboundEmail, error) {
	email := &models.InboundEmail{}
	err := queueDatabaseConnection.
		Where("email_id = ?", emailId).
		First(email).Error

	if err != nil {
		return nil, err
	}
	return email, nil
}

func Worker(certs *tls.Config, entry *queue.Entry, queueDatabaseConnection *gorm.DB) int8 {

	return 2
}
