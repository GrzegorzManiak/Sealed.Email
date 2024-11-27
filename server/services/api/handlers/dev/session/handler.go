package devSession

import (
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/session"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func handler(ctx *gin.Context, session *session.Claims, databaseConnection *gorm.DB) (*Output, helpers.AppError) {
	return &Output{
		ContentSessionID: session.Content.SessionID,
		ContentUserID:    session.Content.UserID,

		Token:     session.Token,
		Signature: session.Signature,

		HeaderVersion: session.Header.Ver,
		HeaderGroup:   session.Header.Group,
		HeaderEat:     session.Header.ExpiresAt,
		HeaderCat:     session.Header.CreatedAt,
		HeaderRat:     session.Header.RefreshAt,
	}, nil
}
