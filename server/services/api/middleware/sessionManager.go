package middleware

import (
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/session"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func SessionManagerMiddleware(ctx *gin.Context, filter *session.APIConfiguration, databaseConnection *gorm.DB) (*session.Claims, helpers.AppError) {
	if filter.Bypass == true {
		filter.SessionRequired = false
		return nil, nil
	}

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

		zap.L().Debug("SessionManagerMiddleware A", zap.Error(err), zap.String("cookie", cookie), zap.Any("filter", filter))
		return nil, helpers.NewNoAccessError("")
	}

	sessionClaims, err := session.ParseSessionToken(cookie)
	if err != nil {
		zap.L().Debug("SessionManagerMiddleware B", zap.Error(err), zap.String("cookie", cookie), zap.Any("filter", filter))
		return nil, helpers.NewNoAccessError("")
	}

	if !sessionClaims.Filter(filter) {
		zap.L().Debug("SessionManagerMiddleware C", zap.Any("sessionClaims", sessionClaims), zap.Any("filter", filter))
		return nil, helpers.NewNoAccessError("")
	}

	if sessionClaims.IsExpired() {
		if filter.SessionRequired == false {
			return nil, nil
		}

		zap.L().Debug("SessionManagerMiddleware D", zap.Any("sessionClaims", sessionClaims), zap.Any("filter", filter))
		return nil, helpers.NewNoAccessError("")
	}

	if sessionClaims.NeedsRefresh() {
		newSessionClaims, err := session.RefreshSessionToken(sessionClaims, databaseConnection)
		if err != nil {
			zap.L().Debug("SessionManagerMiddleware E", zap.Error(err), zap.Any("sessionClaims", sessionClaims), zap.Any("filter", filter))
			return nil, helpers.NewNoAccessError("")
		}

		session.SetSessionCookie(ctx, newSessionClaims)
		return &newSessionClaims, nil
	}

	return &sessionClaims, nil
}
