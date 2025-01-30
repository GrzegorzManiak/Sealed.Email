package domainList

import (
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/services"
)

func Handler(input *Input, data *services.Handler) (*Output, helpers.AppError) {
	domains, total, err := fetchDomainsByUserID(data.User, *input, data.DatabaseConnection)
	if err != nil {
		return nil, err
	}

	return &Output{
		Domains: *parseDomainList(domains),
		Total:   total,
	}, nil
}
