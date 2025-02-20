package domainList

import (
	"github.com/GrzegorzManiak/NoiseBackend/internal/errors"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/services"
)

func Handler(input *Input, data *services.Handler) (*Output, errors.AppError) {
	domains, total, err := fetchDomainsByUserID(data.User, *input, data.DatabaseConnection)
	if err != nil {
		return nil, err
	}

	return &Output{
		Domains: *parseDomainList(domains),
		Total:   total,
	}, nil
}
