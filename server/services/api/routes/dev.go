package routes

import (
	"fmt"
	devEncryption "github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/dev/encryption"
	devSession "github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/dev/session"

	"github.com/GrzegorzManiak/NoiseBackend/services/api/services"
	"github.com/gin-gonic/gin"
)

func DevRoutes(router *gin.Engine, baseRoute *services.BaseRoute) {
	fmt.Println("Registering routes for: dev")
	router.GET("/api/dev/session", func(ctx *gin.Context) {
		services.ExecuteRoute[devSession.Input, devSession.Output](ctx, baseRoute, devSession.SessionFilter, devSession.Handler)
	})
	router.GET("/api/dev/encryption", func(ctx *gin.Context) {
		services.ExecuteRoute[devEncryption.Input, devEncryption.Output](ctx, baseRoute, devEncryption.SessionFilter, devEncryption.Handler)
	})
}
