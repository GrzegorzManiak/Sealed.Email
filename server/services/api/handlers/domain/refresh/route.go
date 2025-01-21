package domainVerify

import (
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/GrzegorzManiak/NoiseBackend/internal/service"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func ExecuteRoute(ctx *gin.Context, databaseConnection *gorm.DB, connPool *service.Pools) {
	data, sessionErr := middleware.SessionManagerMiddleware(ctx, SessionFilter, databaseConnection)
	if sessionErr != nil {
		helpers.ErrorResponse(ctx, sessionErr)
		return
	}

	input, err := helpers.ValidateInputData[Input](ctx)
	if err != nil {
		helpers.ErrorResponse(ctx, err)
		return
	}

	output, err := handler(input, ctx, data.Content.UserID, databaseConnection, connPool)
	if err != nil {
		zap.L().Debug("Error handler", zap.Error(err), zap.Any("input", input), zap.Any("output", output))
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
