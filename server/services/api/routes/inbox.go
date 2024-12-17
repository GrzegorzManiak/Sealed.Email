package routes

import (
	"fmt"
	inboxAdd "github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/inbox/add"
	inboxGet "github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/inbox/get"
	inboxList "github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/inbox/list"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InboxRoutes(router *gin.Engine, databaseConnection *gorm.DB) {
	fmt.Println("Registering routes for: login")
	router.POST("/api/inbox/add", func(ctx *gin.Context) {
		inboxAdd.ExecuteRoute(ctx, databaseConnection)
	})
	router.GET("/api/inbox/list", func(ctx *gin.Context) {
		inboxList.ExecuteRoute(ctx, databaseConnection)
	})
	router.GET("/api/inbox/get", func(ctx *gin.Context) {
		inboxGet.ExecuteRoute(ctx, databaseConnection)
	})
}
