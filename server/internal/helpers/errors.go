package helpers

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type AppError interface {
	Message() string
	Title() string
	Fault() string
	Code() int

	Error() string
}

type BaseError struct {
	message string
	title   string
	fault   string
	code    int
}

func (e BaseError) Fault() string {
	return e.fault
}

func (e BaseError) Message() string {
	return e.message
}

func (e BaseError) Title() string {
	return e.title
}

func (e BaseError) Code() int {
	return e.code
}

func (e BaseError) Error() string {
	return fmt.Sprintf("[%s]: %s (%s)", e.title, e.message, e.fault)
}

func NewBaseError(message string, title string, fault string, code int) BaseError {
	return BaseError{
		message: message,
		title:   title,
		fault:   fault,
		code:    code,
	}
}

func NewUserError(message string, title string) BaseError {
	return NewBaseError(message, title, "user", 400)
}
func NewNoAccessError(message string) BaseError {
	return NewBaseError(message, "You are not allowed to access this resource", "access", 401)
}

func NewServerError(message string, title string) BaseError {
	return NewBaseError(message, title, "server", 500)
}

func DataValidationError(message string) BaseError {
	return NewBaseError(message, "Data validation error", "data", 400)
}

func ErrorResponse(ctx *gin.Context, err AppError) {
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
