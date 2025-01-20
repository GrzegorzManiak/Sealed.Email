package api

import (
	"context"
	"fmt"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	PrimaryDatabase "github.com/GrzegorzManiak/NoiseBackend/database/primary"
	ServiceProvider "github.com/GrzegorzManiak/NoiseBackend/internal/service"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/middleware"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/routes"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/services"
	"github.com/gin-contrib/pprof"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

func Start() {
	zap.L().Info("Starting API service")

	router := gin.Default()
	router.Use(ginzap.Ginzap(zap.L(), time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(zap.L(), true))
	router.Use(middleware.URLCleanerMiddleware())
	router.Use(gin.Recovery())
	pprof.Register(router, "debug/")

	databaseConnection := PrimaryDatabase.InitiateConnection()
	routes.RegisterRoutes(router, databaseConnection)
	routes.LoginRoutes(router, databaseConnection)
	routes.DomainRoutes(router, databaseConnection)
	routes.InboxRoutes(router, databaseConnection)
	routes.DevRoutes(router, databaseConnection)

	serviceUUID, err := uuid.NewUUID()
	if err != nil {
		zap.L().Panic("failed to generate service UUID", zap.Error(err))
	}

	serviceAnnouncement := ServiceProvider.Announcement{
		Id:      serviceUUID.String(),
		Port:    config.Server.Port,
		Host:    config.Server.Host,
		Service: config.Etcd.API,
	}

	etcdContext := context.Background()
	_, err = ServiceProvider.NewEtcdService(etcdContext, &config.Etcd.API, &serviceAnnouncement)
	if err != nil {
		zap.L().Panic("failed to create etcd service", zap.Error(err))
	}

	ServiceProvider.KeepConnectionPoolsAlive(etcdContext, config.Etcd.API)
	ServiceProvider.RegisterCallback("filler", services.PoolCallback)

	err = router.Run(fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port))
	if err != nil {
		zap.L().Panic("failed to run API service", zap.Error(err))
	}
}
