package routes

import (
	"fmt"
	domainAdd "github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/domain/add"
	domainDelete "github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/domain/delete"
	domainModify "github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/domain/modify"
	domainRefresh "github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/domain/refresh"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func DomainRoutes(router *gin.Engine, databaseConnection *gorm.DB) {
	fmt.Println("Registering routes for: login")
	router.POST("/domain/add", func(ctx *gin.Context) {
		domainAdd.ExecuteRoute(ctx, databaseConnection)
	})
	router.DELETE("/domain/delete", func(ctx *gin.Context) {
		domainDelete.ExecuteRoute(ctx, databaseConnection)
	})
	router.PUT("/domain/modify", func(ctx *gin.Context) {
		domainModify.ExecuteRoute(ctx, databaseConnection)
	})
	router.PUT("/domain/refresh", func(ctx *gin.Context) {
		domainRefresh.ExecuteRoute(ctx, databaseConnection)
	})
}
