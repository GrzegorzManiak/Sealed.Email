package services

import (
	"fmt"
	"github.com/GrzegorzManiak/NoiseBackend/internal/errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateInputData[T any](ctx *gin.Context) (*T, errors.AppError) {
	var input T

	if err := ctx.ShouldBindHeader(&input); err != nil {
		return nil, errors.Validation(fmt.Sprintf("header: %s", err.Error()))
	}

	if err := ctx.ShouldBindQuery(&input); err != nil {
		return nil, errors.Validation(fmt.Sprintf("query: %s", err.Error()))
	}

	if ctx.Request.Method != "GET" && ctx.Request.Method != "DELETE" {
		if err := ctx.ShouldBindJSON(&input); err != nil {
			return nil, errors.Validation(fmt.Sprintf("json: %s", err.Error()))
		}
	}

	if err := validate.Struct(input); err != nil {
		return nil, errors.Validation(fmt.Sprintf("struct: %s", err.Error()))
	}

	return &input, nil
}

func ValidateOutputData[Output any](output *Output) errors.AppError {
	if err := validate.Struct(output); err != nil {
		return errors.Validation(err.Error())
	}
	return nil
}
