package register

import (
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func validateInputData(ctx *gin.Context) (*Input, helpers.AppError) {
	var input Input
	if err := ctx.ShouldBindJSON(&input); err != nil {
		return nil, helpers.DataValidationError(err.Error())
	}

	if err := validate.Struct(&input); err != nil {
		return nil, helpers.DataValidationError(err.Error())
	}

	return &input, nil
}

func validateOutputData(output interface{}) helpers.AppError {
	if err := validate.Struct(output); err != nil {
		return helpers.InternalServerError()
	}

	return nil
}

func ExecuteRoute(ctx *gin.Context) {
	input, err := validateInputData(ctx)
	if err != nil {
		helpers.ErrorResponse(ctx, err)
		return
	}

	output, err := handler(input, ctx)
	if err != nil {
		helpers.ErrorResponse(ctx, err)
		return
	}

	if err := validateOutputData(output); err != nil {
		helpers.ErrorResponse(ctx, err)
		return
	}

	helpers.SuccessResponse(ctx, output)
}
