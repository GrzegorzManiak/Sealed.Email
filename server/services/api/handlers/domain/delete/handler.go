package domainDelete

import (
	"github.com/GrzegorzManiak/NoiseBackend/internal/errors"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/services"
)

func Handler(input *Input, data *services.Handler) (*Output, errors.AppError) {
	err := deleteDomain(data.User.ID, input.DomainID, data.DatabaseConnection)
	if err != nil {
		return nil, err
	}

	return &Output{}, nil
}
