package routes

import (
	"fmt"
	devSession "github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/dev/session"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func DevRoutes(router *gin.Engine, databaseConnection *gorm.DB) {
	fmt.Println("Registering routes for: dev")
	router.GET("/api/dev/session", func(ctx *gin.Context) {
		devSession.ExecuteRoute(ctx, databaseConnection)
	})
}
