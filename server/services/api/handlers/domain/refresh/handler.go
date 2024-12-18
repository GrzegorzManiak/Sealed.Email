package domainVerify

import (
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/services"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func handler(data *Input, ctx *gin.Context, userID uint, databaseConnection *gorm.DB) (*Output, helpers.AppError) {

	domainModel, err := fetchDomainByID(userID, data.DomainID, databaseConnection)
	if err != nil {
		return nil, err
	}

	err = services.AddDomainToVerificationQueue(ctx, domainModel)
	sentVerification := true
	if err != nil {
		zap.L().Warn("failed to send verification request", zap.Error(err))
		sentVerification = false
	}

	return &Output{
		VerificationSent: sentVerification,
	}, nil
}
