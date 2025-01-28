package plain

import (
	"github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"gorm.io/gorm"
	"strings"
	"time"
)

func getDomain(
	user *models.User,
	domainID string,
	databaseConnection *gorm.DB,
) (*models.UserDomain, helpers.AppError) {
	var domain models.UserDomain
	result := databaseConnection.
		Where("p_id = ? AND user_id = ?", domainID, user.ID).
		First(&domain)

	if result.Error != nil {
		return &models.UserDomain{}, helpers.NewNotFoundError(
			"We could not find the domain you are looking for. Please try again.",
			"Domain not found!",
		)
	}

	return &domain, nil
}

func createSmtpBodyHeader(key string, value string) string {
	escapedValue := strings.ReplaceAll(value, "\n", "")
	escapedValue = strings.ReplaceAll(escapedValue, "\r", "")
	escapedValue = strings.ReplaceAll(escapedValue, "\"", "\\\"")
	header := key + ": " + escapedValue

	var foldedHeader strings.Builder
	lineLength := 0
	for i, char := range header {
		if lineLength == 76 {
			foldedHeader.WriteString("\r\n ")
			lineLength = 1
		}
		foldedHeader.WriteRune(char)
		lineLength++

		if i == len(header)-1 {
			foldedHeader.WriteString("\r\n")
		}
	}

	return foldedHeader.String()
}

func createMessageID(domain string) string {
	return "<" + helpers.GeneratePublicId() + "@" + domain + ">"
}

func createBasicSmtpHeaders(from string, to string, cc []string) []string {
	headers := make([]string, 0)

	headers = append(headers, createSmtpBodyHeader("From", from))
	headers = append(headers, createSmtpBodyHeader("To", to))

	for _, ccAddress := range cc {
		headers = append(headers, createSmtpBodyHeader("Cc", ccAddress))
	}

	today := time.Now().Format("Mon, 02 Jan 2006 15:04:05")
	headers = append(headers, createSmtpBodyHeader("Date", today))
	headers = append(headers, createSmtpBodyHeader("Subject", "Test Email"))

	return headers
}

func createSmtpBody(headers []string, body string) string {
	var smtpBody strings.Builder

	for _, header := range headers {
		smtpBody.WriteString(header)
		smtpBody.WriteString("\r\n")
	}

	smtpBody.WriteString("\r\n")
	smtpBody.WriteString(body)

	return smtpBody.String()
}
