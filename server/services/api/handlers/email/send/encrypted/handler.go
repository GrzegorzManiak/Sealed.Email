package encrypted

import (
	"github.com/GrzegorzManiak/NoiseBackend/internal/errors"
	"github.com/GrzegorzManiak/NoiseBackend/internal/validation"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/services"
)

func Handler(input *Input, data *services.Handler) (*Output, errors.AppError) {
	fromProvidedDomain, err := validation.ExtractDomainFromEmail(input.From.EmailHash)
	if err != nil {
		return nil, errors.User(err.Error(), "Invalid from email address")
	}

	fromDomain, appErr := getDomain(data.User, input.DomainID, data.DatabaseConnection)
	if appErr != nil {
		return nil, appErr
	}

	if !fromDomain.Verified {
		return nil, errors.Access("You must verify your domain before sending emails.")
	}

	if !validation.CompareDomains(fromDomain.Domain, fromProvidedDomain) {
		return nil, errors.User("From email domain does not match the domain you have verified.", "Invalid from email address")
	}

	messageId, err := sendEmail(input, data, fromDomain)
	if err != nil {
		return nil, errors.User(err.Error(), "Failed to send email")
	}

	return &Output{
		MessageID: messageId,
	}, nil
}
