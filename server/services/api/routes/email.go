package routes

import (
	"fmt"
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
}
