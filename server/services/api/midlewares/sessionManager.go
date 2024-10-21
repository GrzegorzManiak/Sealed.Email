package midlewares

import (
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/internal"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/session"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SessionManagerMiddleware(ctx *gin.Context, filter *session.GroupFilter, databaseConnection *gorm.DB) (*session.Claims, internal.AppError) {
	logger := internal.GetLogger()

	if filter.Allow == nil {
		filter.Allow = []string{}
	}

	if filter.Block == nil {
		filter.Block = []string{}
	}

	cookie, err := ctx.Cookie(config.Session.CookieName)
	if err != nil {
		if filter.SessionRequired == false {
			return nil, nil
		}

		logger.Printf("SessionManagerMiddleware Cookie: %v", err)
		return nil, internal.GenericError{
			Message: "you are not allowed to access this resource",
			ErrCode: 401,
		}
	}

	sessionClaims, err := session.ParseSessionToken(cookie)
	if err != nil {
		logger.Printf("SessionManagerMiddleware ParseSessionToken: %v", err)
		return nil, internal.GenericError{
			Message: "you are not allowed to access this resource",
			ErrCode: 401,
		}
	}

	if !sessionClaims.Filter(filter) {
		logger.Printf("SessionManagerMiddleware: %v", "Filter")
		return nil, internal.GenericError{
			Message: "you are not allowed to access this resource",
			ErrCode: 401,
		}
	}

	if sessionClaims.IsExpired() {
		logger.Printf("SessionManagerMiddleware: %v", "Expired")
		return nil, internal.GenericError{
			Message: "you are not allowed to access this resource",
			ErrCode: 401,
		}
	}

	if sessionClaims.NeedsRefresh() {
		newSessionClaims, err := session.RefreshSessionToken(sessionClaims, databaseConnection)
		if err != nil {
			logger.Printf("SessionManagerMiddleware RefreshSessionToken: %v", err)
			return nil, internal.GenericError{
				Message: "you are not allowed to access this resource",
				ErrCode: 401,
			}
		}

		session.SetSessionCookie(ctx, newSessionClaims)
		return &newSessionClaims, nil
	}

	return &sessionClaims, nil
}
