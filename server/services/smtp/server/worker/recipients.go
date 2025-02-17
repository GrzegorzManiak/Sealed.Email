package worker

import (
	"fmt"
	primaryModels "github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/database/smtp/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"maps"
	"slices"
)

func buildProcessedRecipientsMap(email *models.InboundEmail) map[string]struct{} {
	processedMap := make(map[string]struct{})
	for _, domain := range email.ProcessedSuccessfully {
		processedMap[domain] = struct{}{}
	}
	return processedMap
}

func extractUniqueRecipientDomains(email *models.InboundEmail) []string {
	domains := make(map[string]struct{})
	processedMap := buildProcessedRecipientsMap(email)

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
		if err != nil || slices.Contains(processedSuccessfully, domain) {
			zap.L().Debug("Failed to extract domain from email", zap.Error(err))
			continue
		}
		batchedRecipients[domain] = append(batchedRecipients[domain], to)
	}

	return batchedRecipients
}
