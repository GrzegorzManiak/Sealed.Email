package plain

import (
	"github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/email"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	smtpService "github.com/GrzegorzManiak/NoiseBackend/proto/smtp"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/services"
	"go.uber.org/zap"
	"gorm.io/gorm"
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

func setHeaders(
	headers *email.Headers,
	input *Input,
	fromDomain *models.UserDomain,
) (string, error) {
	headers.From(input.From)
	headers.ReplyTo(input.From)
	headers.To(input.To)
	headers.Cc(input.Cc)
	headers.Date()
	headers.Subject(input.Subject)
	headers.NoiseSignature(input.Signature, input.Nonce)

	if err := headers.InReplyTo(input.InReplyTo); err != nil {
		return "", err
	}

	if err := headers.References(input.References); err != nil {
		return "", err
	}

	return headers.MessageId(fromDomain.Domain), nil
}

func sendEmail(
	input *Input,
	data *services.Handler,
	fromDomain *models.UserDomain,
) (string, helpers.AppError) {
	cc, bcc := email.CleanRecipients(input.To, input.Cc, input.Bcc)
	recipients := email.CombineRecipients(input.To, cc, bcc)
	headers := &email.Headers{}

	messageId, err := setHeaders(headers, input, fromDomain)
	if err != nil {
		return "", helpers.NewUserError(err.Error(), "Failed to set headers")
	}

	zap.L().Debug("Email headers", zap.Any("headers", headers))
	signedEmail, err := email.SignEmailWithDkim(headers, input.Body, fromDomain.Domain, fromDomain.DKIMPrivateKey)
	if err != nil {
		return "", helpers.NewServerError("Failed to sign email. Please try again later.", "Failed to sign email")
	}

	if err := email.Send(data.Context, data.ConnectionPool, &smtpService.Email{
		From:      input.From.Email,
		To:        recipients,
		Body:      []byte(signedEmail),
		Challenge: fromDomain.TxtChallenge,
		Version:   "1.0",
		MessageId: messageId,
		Encrypted: false,
	}); err != nil {
		return "", helpers.NewServerError("Failed to send email. Please try again later.", "Failed to send email")
	}

	return messageId, nil
}
