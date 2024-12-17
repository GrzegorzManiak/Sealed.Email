package inboxGet

import (
	"github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func handler(data *Input, ctx *gin.Context, databaseConnection *gorm.DB, user *models.User) (*Output, helpers.AppError) {

	inbox, appError := getInbox(user, data.InboxID, databaseConnection)
	if appError != nil {
		return nil, appError
	}

	return &Output{
		EmailHash:            inbox.EmailHash,
		SymmetricRootKey:     inbox.SymmetricRootKey,
		AsymmetricPublicKey:  inbox.AsymmetricPublicKey,
		AsymmetricPrivateKey: inbox.AsymmetricPrivateKey,
		EncryptedEmailName:   inbox.EncryptedEmailName,
	}, nil
}
