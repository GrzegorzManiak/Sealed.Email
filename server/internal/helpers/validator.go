package helpers

import (
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

var validate = validator.New()

var validators = map[string]string{
	"EncodedP256Key":      "required,base64rawurl,gte=42,lte=46",
	"UserID":              "required,base64rawurl,gte=40,lte=46",
	"EncodedEncryptedKey": "required,base64rawurl,gte=94,lte=106",
	"PublicID":            "required,alphanum,gte=42,lte=65",
}

func validateField(fl validator.FieldLevel, rule string) bool {
	return validate.Var(fl.Field().String(), rule) == nil
}

func RegisterCustomValidators() {
	for tag, rule := range validators {
		tag, rule := tag, rule // Prevent loop variable issues
		err := validate.RegisterValidation(tag, func(fl validator.FieldLevel) bool {
			return validateField(fl, rule)
		})
		if err != nil {
			zap.L().Panic("failed to register custom validator", zap.String("validator", tag), zap.Error(err))
		}
	}
}
