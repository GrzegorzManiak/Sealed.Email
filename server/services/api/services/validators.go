package services

import (
	"net/http"

	"github.com/GrzegorzManiak/NoiseBackend/internal/errors"
	"github.com/GrzegorzManiak/NoiseBackend/internal/validation"
	"github.com/gin-gonic/gin"
)

func ValidateInputData[T any](ctx *gin.Context) (*T, errors.AppError) {
	var input T

	if err := ctx.ShouldBindHeader(&input); err != nil {
		return nil, errors.Validation("header: " + err.Error())
	}

	if err := ctx.ShouldBindQuery(&input); err != nil {
		return nil, errors.Validation("query: " + err.Error())
	}

	if ctx.Request.Method != http.MethodGet && ctx.Request.Method != http.MethodDelete {
		if err := ctx.ShouldBindJSON(&input); err != nil {
			return nil, errors.Validation("json: " + err.Error())
		}
	}

	if err := validation.CustomValidator.Struct(input); err != nil {
		return nil, errors.Validation("struct: " + err.Error())
	}

	return &input, nil
}

func ValidateOutputData[Output any](output *Output) errors.AppError {
	if err := validation.CustomValidator.Struct(output); err != nil {
		return errors.Validation(err.Error())
	}

	return nil
}
