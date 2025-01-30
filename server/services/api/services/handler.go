package services

import (
	"github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/GrzegorzManiak/NoiseBackend/internal/service"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/middleware"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/session"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type BaseRoute struct {
	DatabaseConnection *gorm.DB
	ConnectionPool     *service.Pools
}

type Handler struct {
	*BaseRoute
	User        *models.User
	Context     *gin.Context
	UserFetched bool
}

func ExecuteRoute[InputType any, OutputType any](
	ctx *gin.Context,
	baseRoute *BaseRoute,
	sessionFilter *session.APIConfiguration,
	handler func(input *InputType, data *Handler) (*OutputType, helpers.AppError),
) {
	// -- Session management
	sessionClaims, sessionErr := middleware.SessionManagerMiddleware(ctx, sessionFilter, baseRoute.DatabaseConnection)
	if sessionErr != nil {
		helpers.ErrorResponse(ctx, sessionErr)
		return
	}

	// -- Input validation
	input, err := helpers.ValidateBodyData[InputType](ctx)
	if err != nil {
		helpers.ErrorResponse(ctx, err)
		return
	}

	// -- Fetch user (if required)
	userFetched := false
	var user models.User

	if sessionFilter.SessionRequired {
		user, err = sessionClaims.FetchUser(baseRoute.DatabaseConnection)
		if err != nil {
			helpers.ErrorResponse(ctx, err)
			return
		}
		userFetched = true
	}

	// -- Call handler
	zap.L().Debug("Calling handler", zap.Any("input", input), zap.Any("user", user))
	output, handlerErr := handler(input, &Handler{
		BaseRoute:   baseRoute,
		Context:     ctx,
		User:        &user,
		UserFetched: userFetched,
	})

	if handlerErr != nil {
		zap.L().Debug("Error handler", zap.Error(err), zap.Any("input", input), zap.Any("user", user))
		helpers.ErrorResponse(ctx, err)
		return
	}

	// -- Output validation
	if outErr := helpers.ValidateOutputData(output); outErr != nil {
		zap.L().Debug("Error ValidateOutputData", zap.Error(outErr), zap.Any("output", output))
		helpers.ErrorResponse(ctx, err)
		return
	}

	// -- Success response
	helpers.SuccessResponse(ctx, output)
}
