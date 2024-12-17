package routes

import (
	"fmt"
	domainAdd "github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/domain/add"
	domainDelete "github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/domain/delete"
	domainGet "github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/domain/get"
	domainList "github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/domain/list"
	domainModify "github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/domain/modify"
	domainRefresh "github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/domain/refresh"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func DomainRoutes(router *gin.Engine, databaseConnection *gorm.DB) {
	fmt.Println("Registering routes for: login")
	router.POST("/api/domain/add", func(ctx *gin.Context) {
		domainAdd.ExecuteRoute(ctx, databaseConnection)
	})
	router.DELETE("/api/domain/delete", func(ctx *gin.Context) {
		domainDelete.ExecuteRoute(ctx, databaseConnection)
	})
	router.PUT("/api/domain/modify", func(ctx *gin.Context) {
		domainModify.ExecuteRoute(ctx, databaseConnection)
	})
	router.PUT("/api/domain/refresh", func(ctx *gin.Context) {
		domainRefresh.ExecuteRoute(ctx, databaseConnection)
	})
	router.GET("/api/domain/list", func(ctx *gin.Context) {
		domainList.ExecuteRoute(ctx, databaseConnection)
	})
	router.GET("/api/domain/get", func(ctx *gin.Context) {
		domainGet.ExecuteRoute(ctx, databaseConnection)
	})
}
