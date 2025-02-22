package delete

import (
	"github.com/GrzegorzManiak/NoiseBackend/internal/errors"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/services"
)

func Handler(input *Input, data *services.Handler) (*Output, errors.AppError) {
	return &Output{}, deleteEmails(input, data)
}
