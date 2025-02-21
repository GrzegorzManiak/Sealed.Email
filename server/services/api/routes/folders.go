package routes

import (
	"fmt"
	folderCreate "github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/folder/create"
	folderDelete "github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/folder/delete"
	folderList "github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/folder/list"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/services"
	"github.com/gin-gonic/gin"
)

func FolderRoutes(router *gin.Engine, baseRoute *services.BaseRoute) {
	fmt.Println("Registering routes for: login")
	router.POST("/api/folder/create", func(ctx *gin.Context) {
		services.ExecuteRoute[folderCreate.Input, folderCreate.Output](ctx, baseRoute, folderCreate.SessionFilter, folderCreate.Handler)
	})
	router.DELETE("/api/folder/delete", func(ctx *gin.Context) {
		services.ExecuteRoute[folderDelete.Input, folderDelete.Output](ctx, baseRoute, folderDelete.SessionFilter, folderDelete.Handler)
	})
	router.GET("/api/folder/list", func(ctx *gin.Context) {
		services.ExecuteRoute[folderList.Input, folderList.Output](ctx, baseRoute, folderList.SessionFilter, folderList.Handler)
	})
}
