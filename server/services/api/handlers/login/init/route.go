package loginInit

import (
	"github.com/GrzegorzManiak/NoiseBackend/internal"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/midlewares"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ExecuteRoute(ctx *gin.Context, databaseConnection *gorm.DB) {
	_, sessionErr := midlewares.SessionManagerMiddleware(ctx, SessionFilter, databaseConnection)
	if sessionErr != nil {
		internal.ErrorResponse(ctx, sessionErr)
		return
	}

	input, err := internal.ValidateInputData[Input](ctx)
	if err != nil {
		internal.ErrorResponse(ctx, err)
		return
	}

	logger := internal.GetLogger()
	output, err := handler(input, ctx, logger, databaseConnection)
	if err != nil {
		logger.Printf("Error handler: %v", err)
		internal.ErrorResponse(ctx, err)
		return
	}

	if err := internal.ValidateOutputData(output); err != nil {
		logger.Printf("Error ValidateOutputData: %v", err)
		internal.ErrorResponse(ctx, err)
		return
	}

	internal.SuccessResponse(ctx, output)
}
