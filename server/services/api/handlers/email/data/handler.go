package data

import (
	"fmt"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/services"
)

func Handler(input *Input, data *services.Handler) (*Output, helpers.AppError) {
	if valid := validateAccessKey(input); !valid {
		return nil, helpers.NewNoAccessError("Invalid access key")
	}

	content := fmt.Sprintf("attachment; filename=%s.eml", input.BucketPath)
	data.Context.Writer.Header().Set("Content-Type", "message/rfc822")
	data.Context.Writer.Header().Set("Content-Disposition", content)
	appErr := fetchEmailData(input, data.MinioClient, &data.Context.Writer)

	return &Output{}, appErr
}
