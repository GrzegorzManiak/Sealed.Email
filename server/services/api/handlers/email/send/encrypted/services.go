package encrypted

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
	from email.EncryptedInbox,
	to email.EncryptedInbox,
	cc []email.EncryptedInbox,
) (string, error) {
	headers.From(input.From.BasicInbox())
	headers.ReplyTo(input.From.BasicInbox())
	headers.To(input.To.BasicInbox())
	headers.Cc(email.ReMapEncryptedInboxes(input.Cc))
	headers.Date()
	headers.EncryptedNoiseSignature(input.Signature, append([]email.EncryptedInbox{from, to}, cc...))
	headers.Subject(input.Subject)

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
	cc, bcc := email.CleanEncryptedRecipients(input.To, input.Cc, input.Bcc)
	recipients := email.FormatEncryptedRecipients(input.To, cc, bcc)
	headers := &email.Headers{}

	messageId, err := setHeaders(headers, input, fromDomain, input.From, input.To, cc)
	if err != nil {
		return "", helpers.NewUserError(err.Error(), "Failed to set headers")
	}

	zap.L().Debug("Email headers", zap.Any("headers", headers))
	signedEmail, err := email.SignEmailWithDkim(headers, input.Body, fromDomain.Domain, fromDomain.DKIMPrivateKey)
	if err != nil {
		return "", helpers.NewServerError("Failed to sign email. Please try again later.", "Failed to sign email")
	}

	if err := email.Send(data.Context, data.ConnectionPool, &smtpService.Email{
		From:      helpers.NormalizeEmail(input.From.BasicInbox().Email),
		To:        recipients,
		InboxKeys: email.ConvertToInboxKeys(bcc),
		Body:      []byte(signedEmail),
		Challenge: fromDomain.TxtChallenge,
		Version:   "1.0",
		MessageId: messageId,
		Encrypted: true,
	}); err != nil {
		return "", helpers.NewServerError("Failed to send email. Please try again later.", "Failed to send email")
	}

	return messageId, nil
}
