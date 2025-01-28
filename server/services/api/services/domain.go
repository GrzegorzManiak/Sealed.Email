package services

import (
	"context"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/GrzegorzManiak/NoiseBackend/internal/service"
	domainService "github.com/GrzegorzManiak/NoiseBackend/proto/domain"
	"go.uber.org/zap"
)

func AddDomainToVerificationQueue(ctx context.Context, connPool *service.Pools, domainModel *models.UserDomain) helpers.AppError {
	pool, err := connPool.GetPool(config.Etcd.Domain.Prefix)
	if err != nil {
		zap.L().Debug("Failed to get domain pool", zap.Error(err))
		return helpers.NewServerError("Sorry! We are unable to process your request at the moment. Please try again later.", "Failed to get domain pool!")
	}

	conn, err := pool.GetConnection()
	if err != nil {
		zap.L().Debug("Failed to get domain client", zap.Error(err))
		return helpers.NewServerError("Sorry! We are unable to process your request at the moment. Please try again later.", "Failed to get domain client!")
	}

	stub := domainService.NewDomainServiceClient(conn.Conn)
	sent, err := stub.QueueDNSVerification(ctx, &domainService.QueueDNSVerificationRequest{
		DomainName:          domainModel.Domain,
		Importance:          10,
		TenantId:            uint64(domainModel.UserID),
		TenantType:          "user",
		TxtVerificationCode: domainModel.TxtChallenge,
		DomainID:            uint64(domainModel.ID),
	})

	if err != nil {
		zap.L().Debug("Failed to queue DNS verification", zap.Error(err))
		conn.Working = false
		return helpers.NewServerError(err.Error(), "Failed to queue DNS verification!")
	}

	if !sent.Acknowledged {
		zap.L().Debug("Failed to queue DNS verification", zap.Error(err))
		conn.Working = false
		return helpers.NewServerError("Sorry! We are unable to process your request at the moment. Please try again later.", "Failed to queue DNS verification!")
	}

	return nil
}
