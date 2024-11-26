package routes

import (
	"fmt"
	loginInit "github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/login/init"
	loginVerify "github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/login/verify"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func LoginRoutes(router *gin.Engine, databaseConnection *gorm.DB) {
	fmt.Println("Registering routes for: login")
	router.PUT("/api/login/init", func(ctx *gin.Context) {
		loginInit.ExecuteRoute(ctx, databaseConnection)
	})
	router.PUT("/api/login/verify", func(ctx *gin.Context) {
		loginVerify.ExecuteRoute(ctx, databaseConnection)
	})
}
