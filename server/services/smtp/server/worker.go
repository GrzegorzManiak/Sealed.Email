package server

import (
	"bytes"
	"context"
	"fmt"
	primaryModels "github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/database/smtp/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/GrzegorzManiak/NoiseBackend/internal/queue"
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"maps"
	"slices"
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

func buildProcessedMap(email *models.InboundEmail) map[string]struct{} {
	processedMap := make(map[string]struct{})
	for _, domain := range email.ProcessedSuccessfully {
		processedMap[domain] = struct{}{}
	}
	return processedMap
}

func extractUniqueRecipientDomains(email *models.InboundEmail) []string {
	domains := make(map[string]struct{})
	processedMap := buildProcessedMap(email)

	for _, recipient := range email.To {
		domain, err := helpers.ExtractDomainFromEmail(recipient)
		if err != nil {
			zap.L().Debug("Failed to extract domain from email", zap.Error(err))
		} else if _, ok := processedMap[domain]; !ok {
			domains[domain] = struct{}{}
		}
	}

	return slices.Collect(maps.Keys(domains))
}

func fetchRecipients(primaryDatabaseConnection *gorm.DB, recipientDomains []string) (*[]primaryModels.UserDomain, error) {
	var fetchedDomains []primaryModels.UserDomain
	if err := primaryDatabaseConnection.
		Where("domain IN ? AND verified = 1", recipientDomains).
		Find(&fetchedDomains).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch recipients: %v", err)
	}
	return &fetchedDomains, nil
}

func batchRecipientsByDomain(tos []string, processedSuccessfully []string) map[string][]string {
	batchedRecipients := make(map[string][]string)

	for _, to := range tos {
		domain, err := helpers.ExtractDomainFromEmail(to)
		if err != nil {
			zap.L().Debug("Failed to extract domain from email", zap.Error(err))
			continue
		}

		if slices.Contains(processedSuccessfully, domain) {
			continue
		}

		if _, ok := batchedRecipients[domain]; !ok {
			batchedRecipients[domain] = make([]string, 0)
		}

		batchedRecipients[domain] = append(batchedRecipients[domain], to)
	}

	return batchedRecipients
}

func insertIntoDatabase(primaryDatabaseConnection *gorm.DB, email *models.InboundEmail, domains *[]primaryModels.UserDomain, inboxes map[string][]string) ([]string, queue.WorkerResponse) {
	successfulInserts := make([]string, 0, len(inboxes))
	for _, domain := range *domains {
		inbox, ok := inboxes[domain.Domain]
		if !ok {
			zap.L().Warn("No inbox found for domain", zap.String("domain", domain.Domain))
			continue
		}

		inserts := make([]primaryModels.UserEmail, 0, len(inbox))
		for _, recipient := range inbox {
			PID := helpers.GeneratePublicId()
			inserts = append(inserts, primaryModels.UserEmail{
				PID:                 PID,
				UserID:              domain.UserID,
				UserDomainID:        domain.ID,
				To:                  recipient,
				ReceivedAt:          email.ReceivedAt,
				OriginallyEncrypted: email.IsEncrypted,
				BucketPath:          email.RefID,
			})
		}

		if err := primaryDatabaseConnection.Create(&inserts).Error; err != nil {
			zap.L().Warn("Failed to insert emails", zap.Error(err))
			return successfulInserts, queue.Failed
		}

		zap.L().Debug("Inserted emails", zap.Any("emails", inserts))
		successfulInserts = append(successfulInserts, domain.Domain)
	}

	return successfulInserts, queue.Verified
}

func insertIntoBucket(minioClient *minio.Client, email *models.InboundEmail) error {
	_, err := minioClient.PutObject(context.Background(), "emails", email.RefID, bytes.NewReader(email.RawData), int64(len(email.RawData)), minio.PutObjectOptions{
		ContentType: "message/rfc822",
	})
	if err != nil {
		zap.L().Debug("Failed to insert email into bucket", zap.Error(err))
		return err
	}
	return nil
}

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

	zap.L().Debug("Fetched recipients", zap.Any("recipients", recipients))
	batchedRecipients := batchRecipientsByDomain(email.To, email.ProcessedSuccessfully)
	successfulInserts, code := insertIntoDatabase(primaryDatabaseConnection, email, recipients, batchedRecipients)
	zap.L().Debug("Successful inserts", zap.Any("domains", successfulInserts), zap.Error(err))

	if !email.InBucket {
		if err := insertIntoBucket(minioClient, email); err != nil {
			return queue.Failed
		}
		zap.L().Debug("Inserted email into bucket")
		email.InBucket = true
	}

	email.ProcessedSuccessfully = append(email.ProcessedSuccessfully, successfulInserts...)
	if err := queueDatabaseConnection.Save(email).Error; err != nil {
		zap.L().Debug("Failed to save email", zap.Error(err))
		return queue.Failed
	}

	return code
}
