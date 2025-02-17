package services

import (
	"github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/GrzegorzManiak/NoiseBackend/internal/service"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/middleware"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/session"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type BaseRoute struct {
	DatabaseConnection *gorm.DB
	ConnectionPool     *service.Pools
	MinioClient        *minio.Client
}

type Handler struct {
	*BaseRoute
	User        *models.User
	Session     *session.Claims
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
		zap.L().Debug("Error SessionManagerMiddleware", zap.Error(sessionErr))
		helpers.ErrorResponse(ctx, sessionErr)
		return
	}

	// -- Input validation
	input, err := helpers.ValidateInputData[InputType](ctx)
	if err != nil {
		zap.L().Debug("Error ValidateInputData", zap.Error(err), zap.Any("input", input))
		helpers.ErrorResponse(ctx, err)
		return
	}

	// -- Fetch user (if required)
	userFetched := false
	var user models.User

	if sessionFilter.SessionRequired {
		user, err = sessionClaims.FetchUser(baseRoute.DatabaseConnection)
		if err != nil {
			zap.L().Debug("Error FetchUser", zap.Error(err))
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
		Session:     sessionClaims,
	})

	if handlerErr != nil {
		zap.L().Debug("Error handler", zap.Error(err), zap.Any("input", input), zap.Any("user", user))
		helpers.ErrorResponse(ctx, handlerErr)
		return
	}

	if sessionFilter.SelfResponse {
		zap.L().Debug("SelfResponse", zap.Any("output", output))
		return
	}

	// -- Output validation
	if outErr := helpers.ValidateOutputData(output); outErr != nil {
		zap.L().Debug("Error ValidateOutputData", zap.Error(outErr), zap.Any("output", output))
		helpers.ErrorResponse(ctx, outErr)
		return
	}

	// -- Success response
	helpers.SuccessResponse(ctx, output)
}
