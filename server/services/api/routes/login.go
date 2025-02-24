package routes

import (
	"fmt"

	loginInit "github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/login/init"
	loginVerify "github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/login/verify"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/services"
	"github.com/gin-gonic/gin"
)

func LoginRoutes(router *gin.Engine, baseRoute *services.BaseRoute) {
	fmt.Println("Registering routes for: login")
	router.PUT("/api/login/init", func(ctx *gin.Context) {
		services.ExecuteRoute[loginInit.Input, loginInit.Output](ctx, baseRoute, loginInit.SessionFilter, loginInit.Handler)
	})

	router.PUT("/api/login/verify", func(ctx *gin.Context) {
		services.ExecuteRoute[loginVerify.Input, loginVerify.Output](ctx, baseRoute, loginVerify.SessionFilter, loginVerify.Handler)
	})
}
