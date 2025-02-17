package worker

import (
	"bytes"
	"context"
	primaryModels "github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/database/smtp/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/GrzegorzManiak/NoiseBackend/internal/queue"
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func getEmailById(emailId string, queueDatabaseConnection *gorm.DB) (*models.InboundEmail, error) {
	email := &models.InboundEmail{}
	if err := queueDatabaseConnection.Where("ref_id = ?", emailId).First(email).Error; err != nil {
		return nil, err
	}
	return email, nil
}

func prepareInserts(email *models.InboundEmail, domain primaryModels.UserDomain, inbox []string) []primaryModels.UserEmail {
	inserts := make([]primaryModels.UserEmail, 0, len(inbox))
	for _, recipient := range inbox {
		inserts = append(inserts, primaryModels.UserEmail{
			PID:                 helpers.GeneratePublicId(64),
			UserID:              domain.UserID,
			UserDomainID:        domain.ID,
			To:                  recipient,
			ReceivedAt:          email.ReceivedAt,
			OriginallyEncrypted: email.IsEncrypted,
			BucketPath:          email.RefID,
			DomainPID:           domain.PID,
		})
	}
	return inserts
}

func insertIntoDatabase(primaryDatabaseConnection *gorm.DB, email *models.InboundEmail, domains *[]primaryModels.UserDomain, inboxes map[string][]string) ([]string, queue.WorkerResponse) {
	successfulInserts := make([]string, 0, len(inboxes))
	for _, domain := range *domains {
		if inbox, ok := inboxes[domain.Domain]; ok {
			inserts := prepareInserts(email, domain, inbox)
			if err := primaryDatabaseConnection.Create(&inserts).Error; err != nil {
				zap.L().Warn("Failed to insert emails", zap.Error(err))
				return successfulInserts, queue.Failed
			}
			zap.L().Debug("Inserted emails", zap.Any("emails", inserts))
			successfulInserts = append(successfulInserts, domain.Domain)
		} else {
			zap.L().Warn("No inbox found for domain", zap.String("domain", domain.Domain))
		}
	}
	return successfulInserts, queue.Verified
}

func insertIntoBucket(minioClient *minio.Client, email *models.InboundEmail) error {
	if _, err := minioClient.PutObject(context.Background(), "emails", email.RefID, bytes.NewReader(email.RawData), int64(len(email.RawData)), minio.PutObjectOptions{
		ContentType: "message/rfc822",
	}); err != nil {
		zap.L().Debug("Failed to insert email into bucket", zap.Error(err))
		return err
	}
	return nil
}
