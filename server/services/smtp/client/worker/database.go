package worker

import (
	"bytes"
	"context"
	primaryModels "github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/database/smtp/models"
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

func insertIntoBucket(minioClient *minio.Client, email *models.OutboundEmail) error {
	if _, err := minioClient.PutObject(context.Background(), "emails", email.RefID, bytes.NewReader(email.Body), int64(len(email.Body)), minio.PutObjectOptions{
		ContentType:  "message/rfc822",
		UserMetadata: map[string]string{"type": "outbound"},
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
