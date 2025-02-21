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
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
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
	minioClient, err := minio.New(config.Bucket.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.Bucket.AccessKey, config.Bucket.SecretKey, ""),
		Secure: config.Bucket.UseSSL,
	})
	if err != nil {
		zap.L().Panic("failed to create minio client", zap.Error(err))
	}
	zap.L().Info("Minio client created", zap.Any("client", minioClient))

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
	etcdService, err := ServiceProvider.NewEtcdService(etcdContext, &config.Etcd.API, &serviceAnnouncement)
	if err != nil {
		zap.L().Panic("failed to create etcd service", zap.Error(err))
	}

	connPool, err := ServiceProvider.NewPools(etcdContext, etcdService, config.Certificates.API)
	if err != nil {
		zap.L().Panic("failed to create distributed service", zap.Error(err))
	}

	baseRoute := &services.BaseRoute{
		DatabaseConnection: databaseConnection,
		ConnectionPool:     connPool,
		MinioClient:        minioClient,
	}

	routes.RegisterRoutes(router, baseRoute)
	routes.LoginRoutes(router, baseRoute)
	routes.DomainRoutes(router, baseRoute)
	routes.DevRoutes(router, baseRoute)
	routes.EmailRoutes(router, baseRoute)
	routes.FolderRoutes(router, baseRoute)

	err = router.Run(fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port))
	if err != nil {
		zap.L().Panic("failed to run API service", zap.Error(err))
	}
}
