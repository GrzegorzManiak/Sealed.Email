package routes

import (
	"fmt"
	devSession "github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/dev/session"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/services"
	"github.com/gin-gonic/gin"
)

func DevRoutes(router *gin.Engine, baseRoute *services.BaseRoute) {
	fmt.Println("Registering routes for: dev")
	router.GET("/api/dev/session", func(ctx *gin.Context) {
		services.ExecuteRoute[devSession.Input, devSession.Output](ctx, baseRoute, devSession.SessionFilter, devSession.Handler)
	})
}
