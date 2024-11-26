package routes

import (
	"fmt"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/register"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(router *gin.Engine, databaseConnection *gorm.DB) {
	fmt.Println("Registering routes for: register")
	router.POST("/api/register", func(ctx *gin.Context) {
		register.ExecuteRoute(ctx, databaseConnection)
	})
}
