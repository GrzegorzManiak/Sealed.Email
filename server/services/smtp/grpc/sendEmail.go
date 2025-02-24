package grpc

import (
	"context"

	"github.com/GrzegorzManiak/NoiseBackend/proto/smtp"
	"github.com/GrzegorzManiak/NoiseBackend/services/smtp/client"
	"go.uber.org/zap"
)

func (s *Server) SendEmail(ctx context.Context, email *smtp.Email) (*smtp.SendEmailResponse, error) {
	data, err := client.QueueEmail(email, s.QueueDatabaseConnection, s.OutboundQueue)
	zap.L().Debug("Queued email", zap.Any("data", data.EmailId))

	if err != nil {
		return &smtp.SendEmailResponse{Success: false}, err
	}

	return &smtp.SendEmailResponse{Success: true}, nil
}
