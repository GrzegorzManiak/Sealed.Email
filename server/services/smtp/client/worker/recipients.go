package worker

import (
	"github.com/GrzegorzManiak/NoiseBackend/database/smtp/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"go.uber.org/zap"
	"slices"
	"strings"
)

func groupRecipients(email *models.OutboundEmail, sentSuccessfully []string) (map[string][]string, error) {
	groupedRecipients := make(map[string][]string)
	bccRecipients := createBccMap(email)

	for _, recipient := range email.To {
		recipient = strings.ToLower(recipient)
		domain, err := helpers.ExtractDomainFromEmail(recipient)
		if err != nil {
			zap.L().Debug("Failed to extract domain from email", zap.Error(err))
			return nil, err
		}

		// -- BCC recipients are not included in the grouped recipients
		if _, ok := bccRecipients[recipient]; ok {
			continue
		}

		if _, ok := groupedRecipients[domain]; !ok {
			groupedRecipients[domain] = []string{}
		}

		if slices.Contains(sentSuccessfully, recipient) {
			continue
		}

		groupedRecipients[domain] = append(groupedRecipients[domain], recipient)
	}

	return groupedRecipients, nil
}

func createBccMap(email *models.OutboundEmail) map[string]struct{} {
	emailKeys := make(map[string]struct{})
	for _, key := range email.OutboundEmailKeys {
		emailKeys[strings.ToLower(key.EmailHash)] = struct{}{}
	}
	return emailKeys
}
