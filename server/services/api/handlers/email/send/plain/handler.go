package plain

import (
	"github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/GrzegorzManiak/NoiseBackend/internal/service"
	smtpService "github.com/GrzegorzManiak/NoiseBackend/proto/smtp"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func handler(data *Input, ctx *gin.Context, databaseConnection *gorm.DB, connPool *service.Pools, user *models.User) (*Output, helpers.AppError) {

	fromProvidedDomain, err := helpers.ExtractDomainFromEmail(data.From)
	if err != nil {
		return nil, helpers.NewUserError(err.Error(), "Invalid from email address")
	}

	fromDomain, appErr := getDomain(user, data.DomainID, databaseConnection)
	if appErr != nil {
		return nil, appErr
	}

	if fromDomain.Verified == false {
		return nil, helpers.NewNoAccessError("You must verify your domain before sending emails.")
	}

	if !helpers.CompareDomains(fromDomain.Domain, fromProvidedDomain) {
		return nil, helpers.NewUserError("From email domain does not match the domain you have verified.", "Invalid from email address")
	}

	headers := createBasicSmtpHeaders(data.From, data.To, data.Cc)
	messageId := createMessageID(fromDomain.Domain)
	headers = append(headers, createSmtpBodyHeader("Message-ID", messageId))
	body := createSmtpBody(headers, data.Body)

	err = services.SendEmail(ctx, connPool, &smtpService.Email{
		From:          data.From,
		To:            []string{data.To},
		Body:          []byte(body),
		DkimSignature: "",
		Version:       "1.0.0",
		InReplyTo:     "",
		References:    nil,
		MessageId:     messageId,
		Encrypted:     false,
	})

	if err != nil {
		return nil, helpers.NewServerError("Failed to send email. Please try again later.", "Failed to send email")
	}

	return &Output{}, nil
}
