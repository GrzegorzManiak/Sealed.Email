package grpc

import (
	"context"
	"github.com/GrzegorzManiak/NoiseBackend/proto/smtp"
	"github.com/GrzegorzManiak/NoiseBackend/services/smtp/client"
)

func (s *Server) SendEmail(ctx context.Context, email *smtp.Email) (*smtp.SendEmailResponse, error) {
	_, err := client.QueueEmail(email, s.QueueDatabaseConnection, s.OutboundQueue)
	if err != nil {
		return &smtp.SendEmailResponse{Success: false}, err
	}

	return &smtp.SendEmailResponse{Success: true}, nil
}
