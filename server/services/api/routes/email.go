package routes

import (
	"fmt"
	emailData "github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/email/data"
	emailGet "github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/email/get"
	emailsList "github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/email/list"
	emailModify "github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/email/modify"
	sendEncrypted "github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/email/send/encrypted"
	sendPlain "github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/email/send/plain"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/services"
	"github.com/gin-gonic/gin"
)

func EmailRoutes(router *gin.Engine, baseRoute *services.BaseRoute) {
	fmt.Println("Registering routes for: email")
	router.POST("/api/email/send/plain", func(ctx *gin.Context) {
		services.ExecuteRoute[sendPlain.Input, sendPlain.Output](ctx, baseRoute, sendPlain.SessionFilter, sendPlain.Handler)
	})
	router.POST("/api/email/send/encrypted", func(ctx *gin.Context) {
		services.ExecuteRoute[sendEncrypted.Input, sendEncrypted.Output](ctx, baseRoute, sendEncrypted.SessionFilter, sendEncrypted.Handler)
	})
	router.GET("/api/email/list", func(ctx *gin.Context) {
		services.ExecuteRoute[emailsList.Input, emailsList.Output](ctx, baseRoute, emailsList.SessionFilter, emailsList.Handler)
	})
	router.GET("/api/email/get", func(ctx *gin.Context) {
		services.ExecuteRoute[emailGet.Input, emailGet.Output](ctx, baseRoute, emailGet.SessionFilter, emailGet.Handler)
	})
	router.GET("/api/email/data", func(ctx *gin.Context) {
		services.ExecuteRoute[emailData.Input, emailData.Output](ctx, baseRoute, emailData.SessionFilter, emailData.Handler)
	})
	router.PUT("/api/email/modify", func(ctx *gin.Context) {
		services.ExecuteRoute[emailModify.Input, emailModify.Output](ctx, baseRoute, emailModify.SessionFilter, emailModify.Handler)
	})
}
