package routes

import (
	"fmt"
	inboxAdd "github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/inbox/add"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InboxRoutes(router *gin.Engine, databaseConnection *gorm.DB) {
	fmt.Println("Registering routes for: login")
	router.POST("/api/inbox/add", func(ctx *gin.Context) {
		inboxAdd.ExecuteRoute(ctx, databaseConnection)
	})
}
