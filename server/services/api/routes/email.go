package routes

import (
	"fmt"
	"github.com/GrzegorzManiak/NoiseBackend/internal/service"
	plainSend "github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/email/send/plain"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func EmailRoutes(router *gin.Engine, databaseConnection *gorm.DB, connPool *service.Pools) {
	fmt.Println("Registering routes for: email")
	router.POST("/api/email/send/plain", func(ctx *gin.Context) {
		plainSend.ExecuteRoute(ctx, databaseConnection, connPool)
	})
}
