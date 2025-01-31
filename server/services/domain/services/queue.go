package services

import (
	database "github.com/GrzegorzManiak/NoiseBackend/database/domain/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/queue"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var QueueName = "Domain Verification Queue"

func Worker(entry *queue.Entry, primaryDatabaseConnection *gorm.DB) int8 {

	data, err := database.UnmarshalVerificationQueue(entry.Data)
	if err != nil {
		zap.L().Error("Failed to unmarshal verification queue", zap.Error(err))
		return 2
	}

	zap.L().Debug("Processing verification queue", zap.Any("entry", entry))
	if err := VerifyDns(data.DomainName, data.TxtVerification); err != nil {
		zap.L().Error("Failed to delete entry", zap.Error(err))
		return 2
	}

	dbErr := VerifyDomain(data.DomainName, data.TenantID, data.DomainID, primaryDatabaseConnection)
	if dbErr != nil {
		zap.L().Debug("Failed to verify domain", zap.Error(dbErr))
		return 2
	}

	zap.L().Debug("Domain verification successful", zap.Any("data", data))
	return 1
}
