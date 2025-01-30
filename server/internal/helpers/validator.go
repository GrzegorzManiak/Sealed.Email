package helpers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

var validate = validator.New()

func ValidateInputData[T interface{}](ctx *gin.Context) (*T, AppError) {
	var input T

	if err := ctx.ShouldBindHeader(&input); err != nil {
		return nil, DataValidationError(fmt.Sprintf("header: %s", err.Error()))
	}

	if err := ctx.ShouldBindQuery(&input); err != nil {
		return nil, DataValidationError(fmt.Sprintf("query: %s", err.Error()))
	}

	if ctx.Request.Method != "GET" && ctx.Request.Method != "DELETE" {
		if err := ctx.ShouldBindJSON(&input); err != nil {
			return nil, DataValidationError(fmt.Sprintf("json: %s", err.Error()))
		}
	}

	if err := validate.Struct(input); err != nil {
		return nil, DataValidationError(fmt.Sprintf("struct: %s", err.Error()))
	}

	return &input, nil
}

func ValidateOutputData[Output any](output *Output) AppError {
	if err := validate.Struct(output); err != nil {
		return DataValidationError(err.Error())
	}

	return nil
}

var (
	P256KeyValidator    = "required,base64,gte=42,lte=46"
	P256KeyValidatorTag = "P256-B64-Key"

	GenericKeyValidator    = "required,base64,gte=42,lte=46"
	GenericKeyValidatorTag = "Generic-B64-Key"

	EncryptedKeyValidator    = "required,base64,gte=102,lte=106"
	EncryptedKeyValidatorTag = "Encrypted-B64-Key"

	PublicIDValidator    = "required,base64,gte=42,lte=46"
	PublicIDValidatorTag = "PublicID"
)

func p256KeyValidation(fl validator.FieldLevel) bool {
	rules := validate.Var(fl.Field().String(), P256KeyValidator)
	return rules == nil
}

func genericKeyValidation(fl validator.FieldLevel) bool {
	rules := validate.Var(fl.Field().String(), GenericKeyValidator)
	return rules == nil
}

func encryptedKeyValidation(fl validator.FieldLevel) bool {
	rules := validate.Var(fl.Field().String(), EncryptedKeyValidator)
	return rules == nil
}

func randomIDValidation(fl validator.FieldLevel) bool {
	rules := validate.Var(fl.Field().String(), PublicIDValidator)
	return rules == nil
}

func RegisterCustomValidators() {
	err := validate.RegisterValidation(P256KeyValidatorTag, p256KeyValidation)
	if err != nil {
		zap.L().Panic("failed to register custom validator", zap.Error(err))
	}

	err = validate.RegisterValidation(EncryptedKeyValidatorTag, encryptedKeyValidation)
	if err != nil {
		zap.L().Panic("failed to register custom validator", zap.Error(err))
	}

	err = validate.RegisterValidation(GenericKeyValidatorTag, genericKeyValidation)
	if err != nil {
		zap.L().Panic("failed to register custom validator", zap.Error(err))
	}

	err = validate.RegisterValidation(PublicIDValidatorTag, randomIDValidation)
	if err != nil {
		zap.L().Panic("failed to register custom validator", zap.Error(err))
	}
}
