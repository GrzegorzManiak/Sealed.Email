package grpc

import (
	"context"
	"fmt"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	database "github.com/GrzegorzManiak/NoiseBackend/database/domain/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/queue"
	"github.com/GrzegorzManiak/NoiseBackend/internal/validation"
	"github.com/GrzegorzManiak/NoiseBackend/proto/domain"
	"github.com/GrzegorzManiak/NoiseBackend/services/domain/services"
	"go.uber.org/zap"
)

func (s *Server) QueueDNSVerification(ctx context.Context, req *domain.QueueDNSVerificationRequest) (*domain.QueueDNSVerificationResponse, error) {
	cleanDomain, domainErr := validation.TrimDomain(req.DomainName)
	if domainErr != nil {
		return &domain.QueueDNSVerificationResponse{
			Message:      "Invalid domain name",
			Acknowledged: false,
		}, nil
	}

	entry, err := queue.Initiate(config.Domain.Service.MaxRetry, config.Domain.Service.RetryInterval, services.QueueName, database.VerificationQueue{
		DomainName:      cleanDomain,
		DkimPublicKey:   req.DkimPublicKey,
		TxtVerification: req.TxtVerificationCode,
		TenantID:        req.TenantId,
		TenantType:      req.TenantType,
		DomainID:        req.DomainID,
	})

	entry.RefID = fmt.Sprintf("%d:%s:%s", req.TenantId, req.TenantType, cleanDomain)
	zap.L().Debug("Initiating DNS verification request", zap.Any("entry", entry))

	if err != nil {
		zap.L().Error("failed to initiate DNS verification request", zap.Error(err))
		return &domain.QueueDNSVerificationResponse{
			Message:      "Failed to initiate DNS verification request",
			Acknowledged: false,
		}, nil
	}

	zap.L().Debug("Adding entry to queue", zap.Any("entry", entry))
	s.Queue.AddEntry(entry)

	zap.L().Debug("DNS verification request queued", zap.Any("entry", entry))
	return &domain.QueueDNSVerificationResponse{
		Message:      "DNS verification queued",
		Acknowledged: true,
	}, nil
}
