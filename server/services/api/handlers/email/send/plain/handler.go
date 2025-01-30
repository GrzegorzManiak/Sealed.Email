package plain

import (
	"github.com/GrzegorzManiak/NoiseBackend/internal/email"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	smtpService "github.com/GrzegorzManiak/NoiseBackend/proto/smtp"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/services"
)

func Handler(input *Input, data *services.Handler) (*Output, helpers.AppError) {
	fromProvidedDomain, err := helpers.ExtractDomainFromEmail(input.From)
	if err != nil {
		return nil, helpers.NewUserError(err.Error(), "Invalid from email address")
	}

	fromDomain, appErr := getDomain(data.User, input.DomainID, data.DatabaseConnection)
	if appErr != nil {
		return nil, appErr
	}

	if fromDomain.Verified == false {
		return nil, helpers.NewNoAccessError("You must verify your domain before sending emails.")
	}

	if !helpers.CompareDomains(fromDomain.Domain, fromProvidedDomain) {
		return nil, helpers.NewUserError("From email domain does not match the domain you have verified.", "Invalid from email address")
	}

	headers := createBasicSmtpHeaders(input.From, input.To, input.Cc)
	messageId := createMessageID(fromDomain.Domain)
	headers = append(headers, createSmtpBodyHeader("Message-ID", messageId))
	body := createSmtpBody(headers, input.Body)

	err = email.Email(data.Context, data.ConnectionPool, &smtpService.Email{
		From:          input.From,
		To:            []string{input.To},
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
