package inboxGet

import (
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/midlewares"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func ExecuteRoute(ctx *gin.Context, databaseConnection *gorm.DB) {
	sessionClaims, sessionErr := midlewares.SessionManagerMiddleware(ctx, SessionFilter, databaseConnection)
	if sessionErr != nil {
		helpers.ErrorResponse(ctx, sessionErr)
		return
	}

	input, err := helpers.ValidateQueryParams[Input](ctx)
	if err != nil {
		helpers.ErrorResponse(ctx, err)
		return
	}

	user, err := sessionClaims.FetchUser(databaseConnection)
	if err != nil {
		helpers.ErrorResponse(ctx, err)
		return
	}

	output, err := handler(input, ctx, databaseConnection, &user)
	if err != nil {
		zap.L().Debug("Error handler", zap.Error(err), zap.Any("input", input), zap.Any("user", user), zap.Any("output", output))
		helpers.ErrorResponse(ctx, err)
		return
	}

	if err := helpers.ValidateOutputData(output); err != nil {
		zap.L().Debug("Error ValidateOutputData", zap.Error(err), zap.Any("output", output))
		helpers.ErrorResponse(ctx, err)
		return
	}

	helpers.SuccessResponse(ctx, output)
}
