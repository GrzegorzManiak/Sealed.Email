package inboxList

import (
	"github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func handler(data *Input, ctx *gin.Context, databaseConnection *gorm.DB, user *models.User) (*Output, helpers.AppError) {
	inboxes, total, err := fetchInboxesByUserID(user.ID, data.DomainPID, *data.Pagination, databaseConnection)
	if err != nil {
		return nil, err
	}

	return &Output{
		Inboxes: *parseInboxList(inboxes),
		Total:   total,
	}, nil
}
