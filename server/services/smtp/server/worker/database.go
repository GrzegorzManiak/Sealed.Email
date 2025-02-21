package worker

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	primaryModels "github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/database/smtp/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/cryptography"
	emailHelper "github.com/GrzegorzManiak/NoiseBackend/internal/email"
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
			OriginallyEncrypted: email.Encrypted,
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

func insertIntoBucket(minioClient *minio.Client, email *[]byte, refID string) error {
	emailBody := *email
	if _, err := minioClient.PutObject(context.Background(), "emails", refID, bytes.NewReader(emailBody), int64(len(emailBody)), minio.PutObjectOptions{
		ContentType: "message/rfc822",
		UserTags:    map[string]string{"type": "inbound"},
	}); err != nil {
		zap.L().Debug("Failed to insert email into bucket", zap.Error(err))
		return err
	}
	return nil
}

func insertEncrypted(minioClient *minio.Client, email *models.InboundEmail, recipients *[]primaryModels.UserDomain) error {
	key, err := emailHelper.CreateInboxKey()
	if err != nil {
		return fmt.Errorf("failed to create inbox key: %w", err)
	}

	keys := make([]*emailHelper.EncryptionKey, 0, len(*recipients))
	for _, recipient := range *recipients {
		encryptedKey, err := emailHelper.EncryptEmailKey(key, recipient.PublicKey)
		if err != nil {
			return fmt.Errorf("failed to encrypt email key: %w", err)
		}

		keys = append(keys, encryptedKey)
	}

	headers := &emailHelper.Headers{}
	headers.EncryptionKeys(keys)
	headers.ContentType("application/json")

	cipherText, iv, err := cryptography.SymEncrypt(email.RawData, key)
	if err != nil {
		return fmt.Errorf("failed to encrypt email body: %w", err)
	}

	compressedCipher := cryptography.Compress(iv, cipherText)
	encodedCipher := base64.RawStdEncoding.EncodeToString(compressedCipher)
	emailBody := emailHelper.FuseHeadersToBody(*headers, encodedCipher)
	email.Encrypted = true
	emailBytes := []byte(emailBody)

	return insertIntoBucket(minioClient, &emailBytes, email.RefID)
}

func ensureEncryptedBucketInsertion(minioClient *minio.Client, email *models.InboundEmail, recipients *[]primaryModels.UserDomain) error {
	if email.InBucket {
		zap.L().Debug("Email already in bucket")
		return nil
	}

	insertFunc := insertEncrypted
	if email.Encrypted {
		insertFunc = func(minioClient *minio.Client, email *models.InboundEmail, recipients *[]primaryModels.UserDomain) error {
			return insertIntoBucket(minioClient, &email.RawData, email.RefID)
		}
	}

	if err := insertFunc(minioClient, email, recipients); err != nil {
		return fmt.Errorf("failed to insert email into bucket: %w", err)
	}

	zap.L().Debug("Inserted email into bucket")
	email.InBucket = true
	return nil
}
