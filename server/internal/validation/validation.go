package validation

import (
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

var CustomValidator = validator.New()

var validators = map[string]string{
	"EncodedP256Key":      "required,base64rawurl,gte=42,lte=46",
	"UserID":              "required,base64rawurl,gte=40,lte=46",
	"EncodedEncryptedKey": "required,base64rawurl,gte=94,lte=106",
	"PublicID":            "required,alphanum,gte=32,lte=65",
}

func validateField(fl validator.FieldLevel, rule string) bool {
	return CustomValidator.Var(fl.Field().Interface(), rule) == nil
}

func RegisterCustomValidators() {
	for tag, rule := range validators {
		zap.L().Debug("Registering custom validator", zap.String("validator", tag), zap.String("rule", rule))

		err := CustomValidator.RegisterValidation(tag, func(fl validator.FieldLevel) bool {
			return validateField(fl, rule)
		})

		if err != nil {
			zap.L().Panic("failed to register custom validator", zap.String("validator", tag), zap.Error(err))
		}
	}
}
