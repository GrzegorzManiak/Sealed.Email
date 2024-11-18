package domainAdd

import (
	"github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/outsideServices"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
)

func handler(data *Input, ctx *gin.Context, logger *log.Logger, databaseConnection *gorm.DB, user *models.User) (*Output, helpers.AppError) {
	domain, err := trimDomain(data.Domain)
	if err != nil {
		return nil, err
	}

	domainModel, err := insertDomain(user, domain, data.EncryptedRootKey, databaseConnection)
	if err != nil {
		return nil, err
	}

	// -- USER CAN RE-VERIFY, NO NEED TO RETURN ERROR
	err = outsideServices.AddDomainToVerificationQueue(ctx, &domainModel)
	sentVerification := true
	if err != nil {
		logger.Printf("Error adding domain to verification queue: %v", err)
		sentVerification = false
	}

	return &Output{
		DomainID:         domainModel.RID,
		DKIMPublicKey:    domainModel.DKIMPublicKey,
		SentVerification: sentVerification,
		TxtRecord:        domainModel.TxtChallenge,
	}, nil
}
