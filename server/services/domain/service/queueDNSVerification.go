package service

import (
	"context"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	database "github.com/GrzegorzManiak/NoiseBackend/database/domain/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/queue"
	"github.com/GrzegorzManiak/NoiseBackend/proto/domain"
)

func (s *Server) QueueDNSVerification(ctx context.Context, req *domain.QueueDNSVerificationRequest) (*domain.QueueDNSVerificationResponse, error) {
	entry, err := queue.Initiate(config.Domain.Service.RetryMax, config.Domain.Service.RetryInterval, QueueName, database.VerificationQueue{})
	if err != nil {
		return &domain.QueueDNSVerificationResponse{
			Message:      "Failed to initiate DNS verification request",
			Acknowledged: false,
		}, nil
	}

	taskId, err := queue.PushToQueue(s.QueueDatabaseConnection, entry)
	if err != nil {
		return &domain.QueueDNSVerificationResponse{
			Message:      "Failed to queue DNS verification request",
			Acknowledged: false,
		}, nil
	}

	return &domain.QueueDNSVerificationResponse{
		Message:      "DNS verification queued: " + taskId,
		Acknowledged: true,
	}, nil
}
