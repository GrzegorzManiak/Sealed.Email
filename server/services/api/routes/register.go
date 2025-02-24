package routes

import (
	"fmt"

	"github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/register"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/services"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, baseRoute *services.BaseRoute) {
	fmt.Println("Registering routes for: register")
	router.POST("/api/register", func(ctx *gin.Context) {
		services.ExecuteRoute[register.Input, register.Output](ctx, baseRoute, register.SessionFilter, register.Handler)
	})
}
