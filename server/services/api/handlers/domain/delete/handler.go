package domainDelete

import (
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/services"
)

func Handler(input *Input, data *services.Handler) (*Output, helpers.AppError) {
	err := deleteDomain(data.User.ID, input.DomainID, data.DatabaseConnection)
	if err != nil {
		return nil, err
	}

	return &Output{}, nil
}
