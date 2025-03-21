package encrypted

import (
	"github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/email"
	"github.com/GrzegorzManiak/NoiseBackend/internal/errors"
	"github.com/GrzegorzManiak/NoiseBackend/internal/validation"
	smtpService "github.com/GrzegorzManiak/NoiseBackend/proto/smtp"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/services"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func getDomain(
	user *models.User,
	domainID string,
	databaseConnection *gorm.DB,
) (*models.UserDomain, errors.AppError) {
	var domain models.UserDomain
	result := databaseConnection.
		Where("p_id = ? AND user_id = ?", domainID, user.ID).
		First(&domain)

	if result.Error != nil {
		return &models.UserDomain{}, errors.NotFound(
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
	headers.NoiseSignature(input.Signature)
	headers.Subject(input.Subject)
	headers.InboxKeys(append([]email.EncryptedInbox{from, to}, cc...))

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
) (string, errors.AppError) {
	cc, bcc := email.CleanEncryptedRecipients(input.To, input.Cc, input.Bcc)
	recipients := email.FormatEncryptedRecipients(input.To, cc, bcc)
	headers := &email.Headers{}

	messageId, err := setHeaders(headers, input, fromDomain, input.From, input.To, cc)
	if err != nil {
		return "", errors.User(err.Error(), "Failed to set headers")
	}

	zap.L().Debug("Email headers", zap.Any("headers", headers))

	signedEmail, err := email.SignEmailWithDkim(headers, email.FoldEmailBody(input.Body), fromDomain.Domain, fromDomain.DKIMPrivateKey)
	if err != nil {
		return "", errors.User("Failed to sign email. Please try again later.", "Failed to sign email")
	}

	if err := email.Send(data.Context, data.ConnectionPool, &smtpService.Email{
		From:          validation.NormalizeEmail(input.From.BasicInbox().Email),
		To:            recipients,
		InboxKeys:     email.ConvertToInboxKeys(bcc),
		Body:          []byte(signedEmail),
		Challenge:     fromDomain.TxtChallenge,
		Version:       "1.0",
		MessageId:     messageId,
		Encrypted:     true,
		FromUserId:    uint64(fromDomain.UserID),
		FromDomainPID: fromDomain.PID,
		FromDomainId:  uint64(fromDomain.ID),
		PublicKey:     fromDomain.PublicKey,
	}); err != nil {
		return "", errors.User("Failed to send email. Please try again later.", "Failed to send email")
	}

	return messageId, nil
}
