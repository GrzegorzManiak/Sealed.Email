package list

import (
	"github.com/GrzegorzManiak/NoiseBackend/internal/errors"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/services"
)

func Handler(input *Input, data *services.Handler) (*Output, errors.AppError) {
	domains, err := fetchEmails(data.User, *input, data.DatabaseConnection)
	if err != nil {
		return nil, err
	}

	return parseEmailList(domains), nil
}
