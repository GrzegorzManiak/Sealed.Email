package loginVerify

import (
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/midlewares"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ExecuteRoute(ctx *gin.Context, databaseConnection *gorm.DB) {
	_, sessionErr := midlewares.SessionManagerMiddleware(ctx, SessionFilter, databaseConnection)
	if sessionErr != nil {
		helpers.ErrorResponse(ctx, sessionErr)
		return
	}

	input, err := helpers.ValidateInputData[Input](ctx)
	if err != nil {
		helpers.ErrorResponse(ctx, err)
		return
	}

	logger := helpers.GetLogger()
	output, err := handler(input, ctx, logger, databaseConnection)
	if err != nil {
		logger.Printf("Error handler: %v", err)
		helpers.ErrorResponse(ctx, err)
		return
	}

	if err := helpers.ValidateOutputData(output); err != nil {
		logger.Printf("Error ValidateOutputData: %v", err)
		helpers.ErrorResponse(ctx, err)
		return
	}

	helpers.SuccessResponse(ctx, output)
}
