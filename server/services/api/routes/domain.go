package routes

import (
	"fmt"
	domainAdd "github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/domain/add"
	domainDelete "github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/domain/delete"
	domainGet "github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/domain/get"
	domainList "github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/domain/list"
	domainModify "github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/domain/modify"
	domainRefresh "github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/domain/refresh"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/services"
	"github.com/gin-gonic/gin"
)

func DomainRoutes(router *gin.Engine, baseRoute *services.BaseRoute) {
	fmt.Println("Registering routes for: login")
	router.POST("/api/domain/add", func(ctx *gin.Context) {
		services.ExecuteRoute[domainAdd.Input, domainAdd.Output](ctx, baseRoute, domainAdd.SessionFilter, domainAdd.Handler)
	})
	router.DELETE("/api/domain/delete", func(ctx *gin.Context) {
		services.ExecuteRoute[domainDelete.Input, domainDelete.Output](ctx, baseRoute, domainDelete.SessionFilter, domainDelete.Handler)
	})
	router.PUT("/api/domain/modify", func(ctx *gin.Context) {
		services.ExecuteRoute[domainModify.Input, domainModify.Output](ctx, baseRoute, domainModify.SessionFilter, domainModify.Handler)
	})
	router.PUT("/api/domain/refresh", func(ctx *gin.Context) {
		services.ExecuteRoute[domainRefresh.Input, domainRefresh.Output](ctx, baseRoute, domainRefresh.SessionFilter, domainRefresh.Handler)
	})
	router.GET("/api/domain/list", func(ctx *gin.Context) {
		services.ExecuteRoute[domainList.Input, domainList.Output](ctx, baseRoute, domainList.SessionFilter, domainList.Handler)
	})
	router.GET("/api/domain/get", func(ctx *gin.Context) {
		services.ExecuteRoute[domainGet.Input, domainGet.Output](ctx, baseRoute, domainGet.SessionFilter, domainGet.Handler)
	})
}
