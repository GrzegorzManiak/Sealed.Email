package api

import (
	PrimaryDatabase "github.com/GrzegorzManiak/NoiseBackend/database/primary"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/midlewares"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/routes"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"log"
)

func Start() {
	log.Printf("------------------ Starting API Service ------------------")
	router := gin.Default()

	router.Use(midlewares.URLCleanerMiddleware())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	pprof.Register(router, "debug/")

	// COORS to allow every origin (Anon function)
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	databaseConnection := PrimaryDatabase.InitiateConnection()
	routes.RegisterRoutes(router, databaseConnection)
	routes.LoginRoutes(router, databaseConnection)
	routes.DomainRoutes(router, databaseConnection)

	err := router.Run()
	if err != nil {
		panic(err)
	}
}
