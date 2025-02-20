package get

import (
	"github.com/GrzegorzManiak/NoiseBackend/internal/errors"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/email/list"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/services"
)

func Handler(input *Input, data *services.Handler) (*Output, errors.AppError) {
	email, err := fetchEmail(data.User, *input, data.DatabaseConnection)
	if err != nil {
		return nil, err
	}

	parsedEmail := list.ParseEmail(email)
	if parsedEmail == nil {
		return nil, errors.Server("Sorry, we could not parse the email.", "Failed to parse email")
	}

	return &Output{
		Email: parsedEmail,
	}, nil
}
