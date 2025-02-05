package client

import (
	"crypto/tls"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/database/smtp/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/GrzegorzManiak/NoiseBackend/internal/queue"
	"github.com/GrzegorzManiak/NoiseBackend/services/domain/services"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"slices"
	"strings"
)

func getEmailById(emailId string, queueDatabaseConnection *gorm.DB) (*models.OutboundEmail, error) {
	email := &models.OutboundEmail{}
	err := queueDatabaseConnection.Where("email_id = ?", emailId).First(email).Error
	if err != nil {
		return nil, err
	}
	return email, nil
}

func groupRecipients(email *models.OutboundEmail, sentSuccessfully []string, bccRecipients map[string]models.OutboundEmailKeys) (map[string][]string, error) {
	groupedRecipients := make(map[string][]string)

	for _, recipient := range email.To {
		recipient = strings.ToLower(recipient)
		domain, err := helpers.ExtractDomainFromEmail(recipient)
		if err != nil {
			zap.L().Debug("Failed to extract domain from email", zap.Error(err))
			return nil, err
		}

		if _, ok := bccRecipients[domain]; !ok {
			continue
		}

		if _, ok := groupedRecipients[domain]; !ok {
			groupedRecipients[domain] = []string{}
		}

		if slices.Contains(sentSuccessfully, recipient) {
			continue
		}

		groupedRecipients[domain] = append(groupedRecipients[domain], recipient)
	}

	return groupedRecipients, nil
}

func batchSendEmails(certs *tls.Config, email *models.OutboundEmail, domain string, recipients []string) error {
	var batch []string
	for i, recipient := range recipients {
		batch = append(batch, recipient)
		if config.Smtp.MaxOutboundRecipients == i+1 || i+1 == len(recipients) {
			zap.L().Debug("Sending batch of emails", zap.Any("batch", batch), zap.String("domain", domain))
			if err := attemptSendEmail(certs, email, domain, batch); err != nil {
				zap.L().Debug("Failed to send email", zap.Error(err))
				return err
			}
			batch = []string{}
		}
	}
	return nil
}

func createBccMap(email *models.OutboundEmail) map[string]models.OutboundEmailKeys {
	emailKeys := make(map[string]models.OutboundEmailKeys)
	for _, key := range email.OutboundEmailKeys {
		emailKeys[key.EmailHash] = key
	}
	return emailKeys
}

func sendEmails(certs *tls.Config, email *models.OutboundEmail, groupedRecipients map[string][]string) (int8, []string) {
	var sentSuccessfully []string
	for domain, recipients := range groupedRecipients {
		if slices.Contains(sentSuccessfully, domain) {
			continue
		}
		zap.L().Debug("Sending email to domain", zap.String("domain", domain), zap.Any("recipients", recipients))
		if err := batchSendEmails(certs, email, domain, recipients); err != nil {
			zap.L().Debug("Failed to batch send emails", zap.Error(err))
			return 2, sentSuccessfully
		} else {
			zap.L().Debug("Batch sent successfully")
			sentSuccessfully = append(sentSuccessfully, domain)
		}
	}

	//
	// This may look like duplicated code, but it's not. The previous loop sends emails to domains,
	// this one sends emails to encrypted bcc recipients, which cant be batched as each email needs
	// to be sent with an associated key, which would expose the fact that the email is bcc'd.
	//
	for _, bccKeys := range email.OutboundEmailKeys {
		if slices.Contains(sentSuccessfully, bccKeys.EmailHash) {
			continue
		}
		zap.L().Debug("Sending email to bcc", zap.Any("bccKeys", bccKeys))
		if err := attemptSendEmailBcc(certs, email, bccKeys.EmailHash, bccKeys); err != nil {
			zap.L().Debug("Failed to send email to bcc", zap.Error(err))
			return 2, sentSuccessfully
		} else {
			zap.L().Debug("Bcc sent successfully")
			sentSuccessfully = append(sentSuccessfully, bccKeys.EmailHash)
		}
	}

	return 1, sentSuccessfully
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

	bccRecipients := createBccMap(email)
	groupedRecipients, err := groupRecipients(email, email.SentSuccessfully, bccRecipients)
	zap.L().Debug("Grouped recipients", zap.Any("groupedRecipients", groupedRecipients))
	if err != nil {
		zap.L().Debug("Failed to group recipients", zap.Error(err))
		return 2
	}

	fromDomain, err := helpers.ExtractDomainFromEmail(email.From)
	if err != nil {
		zap.L().Debug("Failed to extract domain from email", zap.Error(err))
		return 2
	}

	if err = services.VerifyDns(fromDomain, email.Challenge); err != nil {
		zap.L().Debug("Failed to verify dns", zap.Error(err))
		return 2
	}

	code, sentSuccessfully := sendEmails(certs, email, groupedRecipients)
	email.SentSuccessfully = sentSuccessfully
	if err := queueDatabaseConnection.Save(email).Error; err != nil {
		zap.L().Debug("Failed to save email", zap.Error(err))
		return 2
	}

	zap.L().Debug("Email sent", zap.Any("email id", emailId), zap.Any("recipients", sentSuccessfully), zap.Any("code", code))
	return code
}
