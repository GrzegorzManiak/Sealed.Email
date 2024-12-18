package devSession

import (
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/midlewares"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func ExecuteRoute(ctx *gin.Context, databaseConnection *gorm.DB) {
	data, sessionErr := midleware.SessionManagerMiddleware(ctx, SessionFilter, databaseConnection)
	if sessionErr != nil {
		helpers.ErrorResponse(ctx, sessionErr)
		return
	}

	output, err := handler(ctx, data, databaseConnection)
	if err != nil {
		zap.L().Debug("Error handler", zap.Error(err), zap.Any("output", output))
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
