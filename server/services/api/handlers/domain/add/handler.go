package domainAdd

import (
	"github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
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

	return &Output{
		Verification:  domainModel.RID,
		DKIMPublicKey: domainModel.DKIMPublicKey,
	}, nil
}
