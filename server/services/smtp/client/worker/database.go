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
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func getEmailById(emailId string, queueDatabaseConnection *gorm.DB) (*models.OutboundEmail, error) {
	email := &models.OutboundEmail{}
	err := queueDatabaseConnection.
		Preload("OutboundEmailKeys").
		Where("email_id = ?", emailId).
		First(email).Error

	if err != nil {
		return nil, err
	}
	return email, nil
}

func insertIntoBucket(minioClient *minio.Client, email *[]byte, refID string) error {
	emailBody := *email
	if _, err := minioClient.PutObject(context.Background(), "emails", refID, bytes.NewReader(emailBody), int64(len(emailBody)), minio.PutObjectOptions{
		ContentType: "message/rfc822",
		UserTags:    map[string]string{"type": "outbound"},
	}); err != nil {
		zap.L().Debug("Failed to insert email into bucket", zap.Error(err))
		return err
	}
	return nil
}

func insertIntoDatabase(primaryDatabaseConnection *gorm.DB, email *models.OutboundEmail) error {
	if len(email.To) == 0 {
		return nil
	}

	insert := primaryModels.UserEmail{
		PID:                 email.RefID,
		UserID:              email.FromUserId,
		UserDomainID:        email.FromDomainId,
		To:                  email.To[0],
		ReceivedAt:          email.CreatedAt.Unix(),
		BucketPath:          email.RefID,
		DomainPID:           email.FromDomainPID,
		OriginallyEncrypted: email.Encrypted,
		Sent:                true,
	}

	if err := primaryDatabaseConnection.Create(&insert).Error; err != nil {
		zap.L().Warn("Failed to insert emails", zap.Error(err))
		return err
	}

	zap.L().Debug("Inserted emails", zap.Any("emails", insert))
	return nil
}

func insertEncrypted(minioClient *minio.Client, email *models.OutboundEmail) error {
	zap.L().Debug("Email is not encrypted, encrypting...")

	key, err := emailHelper.CreateInboxKey()
	if err != nil {
		return fmt.Errorf("failed to create inbox key: %w", err)
	}

	encryptedKey, err := emailHelper.EncryptEmailKey(key, email.PublicKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt email key: %w", err)
	}

	headers := &emailHelper.Headers{}
	headers.EncryptionKeys([]*emailHelper.EncryptionKey{encryptedKey})
	headers.ContentType("application/json")

	cipherText, iv, err := cryptography.SymEncrypt(email.Body, key)
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

func ensureEncryptedBucketInsertion(minioClient *minio.Client, email *models.OutboundEmail) error {
	if email.InBucket {
		zap.L().Debug("Email already in bucket")
		return nil
	}

	insertFunc := insertEncrypted
	if email.Encrypted {
		insertFunc = func(minioClient *minio.Client, email *models.OutboundEmail) error {
			return insertIntoBucket(minioClient, &email.Body, email.RefID)
		}
	}

	if err := insertFunc(minioClient, email); err != nil {
		return fmt.Errorf("failed to insert email into bucket: %w", err)
	}

	zap.L().Debug("Inserted email into bucket")
	email.InBucket = true
	return nil
}
