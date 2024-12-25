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
	dnsRecords, err := FetchDnsRecords(data.DomainName)
	if err != nil {
		zap.L().Debug("Failed to fetch DNS records", zap.Error(err))
		return 2
	}

	if !MatchTxtRecords(data.TxtVerification, dnsRecords) {
		zap.L().Debug("Failed to match TXT records", zap.Any("data", data), zap.Any("dnsRecords", dnsRecords))
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
