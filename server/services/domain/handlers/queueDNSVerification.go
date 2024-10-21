package handlers

import (
	"context"
	"fmt"
	"github.com/GrzegorzManiak/NoiseBackend/proto/domain"
)

func (s *Service) QueueDNSVerification(ctx context.Context, req *domain.QueueDNSVerificationRequest) (*domain.QueueDNSVerificationResponse, error) {
	fmt.Printf("Verifying DNS for domain: %s\n", req.DomainName)
	return &domain.QueueDNSVerificationResponse{
		Message: "DNS verification queued",
	}, nil
}
