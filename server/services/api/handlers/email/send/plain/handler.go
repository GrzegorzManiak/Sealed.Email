package plain

import (
	"github.com/GrzegorzManiak/NoiseBackend/internal/email"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	smtpService "github.com/GrzegorzManiak/NoiseBackend/proto/smtp"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/services"
	"go.uber.org/zap"
)

func Handler(input *Input, data *services.Handler) (*Output, helpers.AppError) {
	fromProvidedDomain, err := helpers.ExtractDomainFromEmail(input.From.Email)
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

	headers := &email.Headers{}
	headers.From(input.From)
	headers.To(input.To)
	headers.Cc(input.Cc)
	headers.Date()
	headers.Subject(input.Subject)
	headers.NoiseSignature(input.Signature, input.Nonce)
	messageId := headers.MessageId(fromDomain.Domain)

	zap.L().Debug("Email headers", zap.Any("headers", headers))
	signedEmail, err := email.SignEmailWithDkim(headers, input.Body, fromDomain.Domain, fromDomain.DKIMPrivateKey)
	if err != nil {
		return nil, helpers.NewServerError("Failed to sign email. Please try again later.", "Failed to sign email")
	}

	to := []string{input.To.Email}
	for _, cc := range input.Cc {
		to = append(to, cc.Email)
	}

	for _, bcc := range input.Bcc {
		to = append(to, bcc.Email)
	}

	err = email.Email(data.Context, data.ConnectionPool, &smtpService.Email{
		From:      input.From.Email,
		To:        to,
		Body:      []byte(signedEmail),
		Version:   "1.0",
		MessageId: messageId,
		Encrypted: false,
	})

	if err != nil {
		return nil, helpers.NewServerError("Failed to send email. Please try again later.", "Failed to send email")
	}

	return &Output{}, nil
}
