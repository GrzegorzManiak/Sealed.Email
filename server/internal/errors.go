package internal

import (
	"github.com/gin-gonic/gin"
)

type AppError interface {
	Error() string
	Code() int
}

type GenericError struct {
	Message string
	ErrCode int
}

func (e GenericError) Error() string {
	return e.Message
}

func (e GenericError) Code() int {
	return e.ErrCode
}

func DataValidationError(field string) AppError {
	message := "Field '" + field + "' is invalid, please change it and try again."

	return GenericError{
		Message: message,
		ErrCode: 400,
	}
}

func InternalServerError() AppError {
	return GenericError{
		Message: "Internal server error, please try again later.",
		ErrCode: 500,
	}
}

func ErrorResponse(ctx *gin.Context, err AppError) {
	ctx.JSON(err.Code(), gin.H{
		"message": err.Error(),
	})
}

func SuccessResponse(ctx *gin.Context, data interface{}) {
	ctx.Header("Content-Type", "application/json")
	ctx.JSON(200, data)
}

func RedirectResponse(ctx *gin.Context, url string) {
	ctx.Redirect(302, url)
}
