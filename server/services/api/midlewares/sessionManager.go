package midlewares

import (
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/session"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func SessionManagerMiddleware(ctx *gin.Context, filter *session.APIConfiguration, databaseConnection *gorm.DB) (*session.Claims, helpers.AppError) {

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

		zap.L().Debug("SessionManagerMiddleware", zap.Error(err), zap.String("cookie", cookie), zap.Any("filter", filter))
		return nil, helpers.GenericError{
			Message: "you are not allowed to access this resource",
			ErrCode: 401,
		}
	}

	sessionClaims, err := session.ParseSessionToken(cookie)
	if err != nil {
		zap.L().Debug("SessionManagerMiddleware", zap.Error(err), zap.String("cookie", cookie), zap.Any("filter", filter))
		return nil, helpers.GenericError{
			Message: "you are not allowed to access this resource",
			ErrCode: 401,
		}
	}

	if !sessionClaims.Filter(filter) {
		zap.L().Debug("SessionManagerMiddleware", zap.Any("sessionClaims", sessionClaims), zap.Any("filter", filter))
		return nil, helpers.GenericError{
			Message: "you are not allowed to access this resource",
			ErrCode: 401,
		}
	}

	if sessionClaims.IsExpired() {
		zap.L().Debug("SessionManagerMiddleware", zap.Any("sessionClaims", sessionClaims), zap.Any("filter", filter))
		return nil, helpers.GenericError{
			Message: "you are not allowed to access this resource",
			ErrCode: 401,
		}
	}

	if sessionClaims.NeedsRefresh() {
		newSessionClaims, err := session.RefreshSessionToken(sessionClaims, databaseConnection)
		if err != nil {
			zap.L().Debug("SessionManagerMiddleware", zap.Error(err), zap.Any("sessionClaims", sessionClaims), zap.Any("filter", filter))
			return nil, helpers.GenericError{
				Message: "you are not allowed to access this resource",
				ErrCode: 401,
			}
		}

		session.SetSessionCookie(ctx, newSessionClaims)
		return &newSessionClaims, nil
	}

	return &sessionClaims, nil
}
