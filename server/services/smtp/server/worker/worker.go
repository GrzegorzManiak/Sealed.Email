package worker

import (
	"github.com/GrzegorzManiak/NoiseBackend/internal/queue"
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Worker(entry *queue.Entry, queueDatabaseConnection *gorm.DB, primaryDatabaseConnection *gorm.DB, minioClient *minio.Client) queue.WorkerResponse {
	email, err := getEmailById(entry.RefID, queueDatabaseConnection)
	if err != nil {
		zap.L().Debug("Failed to get email by id", zap.Error(err))
		return queue.Failed
	}

	recipientDomains := extractUniqueRecipientDomains(email)
	recipients, err := fetchRecipients(primaryDatabaseConnection, recipientDomains)
	if err != nil {
		zap.L().Debug("Failed to fetch recipients", zap.Error(err))
		return queue.Failed
	}

	if err := ensureEncryptedBucketInsertion(minioClient, email, recipients); err != nil {
		return queue.Failed
	}
	batchedRecipients := batchRecipientsByDomain(email.To, email.ProcessedSuccessfully)
	successfulInserts, code := insertIntoDatabase(primaryDatabaseConnection, email, recipients, batchedRecipients)

	email.ProcessedSuccessfully = append(email.ProcessedSuccessfully, successfulInserts...)
	if err := queueDatabaseConnection.Save(email).Error; err != nil {
		zap.L().Debug("Failed to save email", zap.Error(err))
		return queue.Failed
	}

	return code
}
