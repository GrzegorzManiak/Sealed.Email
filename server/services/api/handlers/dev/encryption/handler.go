package devEncryption

import (
	"github.com/GrzegorzManiak/NoiseBackend/internal/errors"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/services"
)

func Handler(input *Input, data *services.Handler) (*Output, errors.AppError) {
	testData, err := GenerateTestData()
	if err != nil {
		return nil, errors.User(err.Error(), "Failed to generate test data")
	}

	return testData, nil
}
