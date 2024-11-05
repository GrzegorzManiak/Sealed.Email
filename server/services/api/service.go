package api

import (
	"context"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	PrimaryDatabase "github.com/GrzegorzManiak/NoiseBackend/database/primary"
	"github.com/GrzegorzManiak/NoiseBackend/internal/services"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/midlewares"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/routes"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	serviceUUID, err := uuid.NewUUID()
	if err != nil {
		log.Fatalf("failed to generate service UUID: %v", err)
	}

	etcdService := services.ServiceAnnouncement{
		Id:      serviceUUID.String(),
		Port:    config.Server.Port,
		Host:    config.Server.Host,
		Service: config.Etcd.API,
	}

	etcdClient := services.GetEtcdClient(config.Etcd.API)
	etcdContext := context.Background()
	services.KeepLeaseAlive(etcdContext, etcdClient, etcdService, false)
	services.KeepConnectionPoolsAlive(etcdContext, etcdClient)

	err = router.Run()
	if err != nil {
		panic(err)
	}
}
