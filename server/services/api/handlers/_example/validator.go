package register

import (
	"github.com/GrzegorzManiak/NoiseBackend/internal"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func validateInputData(ctx *gin.Context) (*Input, internal.AppError) {
	var input Input
	if err := ctx.ShouldBindJSON(&input); err != nil {
		return nil, internal.DataValidationError(err.Error())
	}

	if err := validate.Struct(&input); err != nil {
		return nil, internal.DataValidationError(err.Error())
	}

	return &input, nil
}

func validateOutputData(output interface{}) internal.AppError {
	if err := validate.Struct(output); err != nil {
		return internal.InternalServerError()
	}

	return nil
}

func ExecuteRoute(ctx *gin.Context) {
	input, err := validateInputData(ctx)
	if err != nil {
		internal.ErrorResponse(ctx, err)
		return
	}

	output, err := handler(input, ctx)
	if err != nil {
		internal.ErrorResponse(ctx, err)
		return
	}

	if err := validateOutputData(output); err != nil {
		internal.ErrorResponse(ctx, err)
		return
	}

	internal.SuccessResponse(ctx, output)
}
