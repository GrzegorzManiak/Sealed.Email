package worker

import (
	"crypto/tls"

	"github.com/GrzegorzManiak/NoiseBackend/database/smtp/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/queue"
	"github.com/GrzegorzManiak/NoiseBackend/internal/validation"
	"github.com/GrzegorzManiak/NoiseBackend/services/domain/services"
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func saveQueueEmail(queueDatabaseConnection *gorm.DB, email *models.OutboundEmail) error {
	if err := queueDatabaseConnection.Save(email).Error; err != nil {
		zap.L().Debug("Failed to save email", zap.Error(err))

		return err
	}

	return nil
}

func Worker(certs *tls.Config, entry *queue.Entry, queueDatabaseConnection *gorm.DB, primaryDatabaseConnection *gorm.DB, minioClient *minio.Client) queue.WorkerResponse {
	zap.L().Debug("Processing smtp queue", zap.Any("entry", entry))

	emailId, err := models.UnmarshalQueueEmailId(entry.Data)
	if err != nil {
		zap.L().Debug("Failed to unmarshal email id", zap.Error(err))

		return queue.Failed
	}

	email, err := getEmailById(emailId.EmailId, queueDatabaseConnection)
	if err != nil {
		zap.L().Debug("Failed to get email by id", zap.Error(err))

		return queue.Failed
	}

	groupedRecipients, err := groupRecipients(email, email.SentSuccessfully)
	if err != nil {
		zap.L().Debug("Failed to group recipients", zap.Error(err))

		return queue.Failed
	}

	fromDomain, err := validation.ExtractDomainFromEmail(email.From)
	if err != nil {
		zap.L().Debug("Failed to extract domain from email", zap.Error(err))

		return queue.Failed
	}

	if err = services.VerifyDns(fromDomain, email.Challenge); err != nil {
		zap.L().Debug("Failed to verify dns", zap.Error(err))

		return queue.Failed
	}

	if err = ensureEncryptedBucketInsertion(minioClient, email); err != nil {
		zap.L().Debug("Failed to insert email into bucket", zap.Error(err))

		return queue.Failed
	}

	if !email.InDatabase {
		if err := insertIntoDatabase(primaryDatabaseConnection, email); err != nil {
			err = saveQueueEmail(queueDatabaseConnection, email)
			zap.L().Debug("Attempt to save email", zap.Error(err))

			return queue.Failed
		}

		zap.L().Debug("Inserted email into database")

		email.InDatabase = true
	}

	code, sentSuccessfully := sendEmails(certs, email, groupedRecipients)

	email.SentSuccessfully = sentSuccessfully
	if err := saveQueueEmail(queueDatabaseConnection, email); err != nil {
		zap.L().Debug("Failed to save email", zap.Error(err))

		return queue.Failed
	}

	zap.L().Debug("Email sent", zap.Any("email id", emailId), zap.Any("recipients", sentSuccessfully), zap.Any("code", code))

	return code
}
