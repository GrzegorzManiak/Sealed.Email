package domainVerify

import (
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/services"
	"go.uber.org/zap"
)

func Handler(input *Input, data *services.Handler) (*Output, helpers.AppError) {

	domainModel, err := fetchDomainByID(data.User.ID, input.DomainID, data.DatabaseConnection)
	if err != nil {
		return nil, err
	}

	err = services.AddDomainToVerificationQueue(data.Context, data.ConnectionPool, domainModel)
	sentVerification := true
	if err != nil {
		zap.L().Warn("failed to send verification request", zap.Error(err))
		sentVerification = false
	}

	return &Output{
		VerificationSent: sentVerification,
	}, nil
}
