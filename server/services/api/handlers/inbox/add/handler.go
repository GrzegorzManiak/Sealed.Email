package inboxAdd

import (
	"github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func handler(data *Input, ctx *gin.Context, databaseConnection *gorm.DB, user *models.User) (*Output, helpers.AppError) {

	domain, appError := getDomain(user, data.DomainID, databaseConnection)
	if appError != nil {
		return nil, appError
	}

	inbox, appError := createInbox(user, domain, data, databaseConnection)
	if appError != nil {
		return nil, appError
	}

	return &Output{
		InboxID: inbox.PID,
	}, nil
}
