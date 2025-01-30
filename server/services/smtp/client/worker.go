package client

import (
	"crypto/tls"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/database/smtp/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
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

func GroupRecipients(email *models.OutboundEmail) (map[string][]string, error) {
	groupedRecipients := make(map[string][]string)
	for _, recipient := range email.To {
		domain, err := helpers.ExtractDomainFromEmail(recipient)
		if err != nil {
			zap.L().Debug("Failed to extract domain from email", zap.Error(err))
			return nil, err
		}

		if _, ok := groupedRecipients[domain]; !ok {
			groupedRecipients[domain] = []string{}
		}

		groupedRecipients[domain] = append(groupedRecipients[domain], recipient)
	}
	return groupedRecipients, nil
}

func BatchSendEmails(certs *tls.Config, email *models.OutboundEmail, domain string, recipients []string) error {
	var batch []string
	for i, recipient := range recipients {
		batch = append(batch, recipient)
		if config.Smtp.MaxOutboundRecipients == i+1 || i+1 == len(recipients) {
			zap.L().Debug("Sending batch of emails", zap.Any("batch", batch), zap.String("domain", domain))
			if err := attemptSendEmail(certs, email, domain); err != nil {
				zap.L().Debug("Failed to send email", zap.Error(err))
				return err
			}
			batch = []string{}
		}
	}
	return nil
}

func Worker(certs *tls.Config, entry *queue.Entry, queueDatabaseConnection *gorm.DB) int8 {
	zap.L().Debug("Processing smtp queue", zap.Any("entry", entry))

	emailId, err := models.UnmarshalQueueEmailId(entry.Data)
	if err != nil {
		zap.L().Debug("Failed to unmarshal email id", zap.Error(err))
		return 2
	}
	zap.L().Debug("Unmarshalled email id", zap.Any("emailId", emailId))

	email, err := getEmailById(emailId.EmailId, queueDatabaseConnection)
	if err != nil {
		zap.L().Debug("Failed to get email by id", zap.Error(err))
		return 2
	}
	zap.L().Debug("Got email by id", zap.Any("email", email))

	groupedRecipients, err := GroupRecipients(email)
	if err != nil {
		zap.L().Debug("Failed to group recipients", zap.Error(err))
		return 2
	}

	for domain, recipients := range groupedRecipients {
		zap.L().Debug("Sending email to domain", zap.String("domain", domain), zap.Any("recipients", recipients))
		if err := BatchSendEmails(certs, email, domain, recipients); err != nil {
			zap.L().Debug("Failed to batch send emails", zap.Error(err))
			return 2
		}
	}

	zap.L().Debug("Email sent", zap.Any("email", email))
	return 1
}
