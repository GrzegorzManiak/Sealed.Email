package server

import (
	"fmt"
	primaryModels "github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/database/smtp/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/GrzegorzManiak/NoiseBackend/internal/queue"
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

func Worker(entry *queue.Entry, queueDatabaseConnection *gorm.DB, primaryDatabaseConnection *gorm.DB) queue.WorkerResponse {
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

	email.ProcessedSuccessfully = append(email.ProcessedSuccessfully, recipientDomains...)
	if err := queueDatabaseConnection.Save(email).Error; err != nil {
		zap.L().Debug("Failed to save email", zap.Error(err))
		return queue.Failed
	}

	return queue.Verified
}
