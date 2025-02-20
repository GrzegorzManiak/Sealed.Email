package helpers

import (
	"github.com/GrzegorzManiak/NoiseBackend/internal/errors"
	"github.com/gin-gonic/gin"
)

func ErrorResponse(ctx *gin.Context, err errors.AppError) {
	ctx.JSON(err.Code(), gin.H{
		"error": true,
		"title": err.Title(),
		"fault": err.Fault(),
		"msg":   err.Message(),
	})
}

func SuccessResponse(ctx *gin.Context, data interface{}) {
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(200, data)
}

func RedirectResponse(ctx *gin.Context, url string) {
	ctx.Redirect(302, url)
}
