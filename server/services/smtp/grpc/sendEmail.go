package grpc

import (
	"context"
	"github.com/GrzegorzManiak/NoiseBackend/proto/smtp"
)

func (s *Server) SendEmail(ctx context.Context, email *smtp.Email) (*smtp.SendEmailResponse, error) {
	return &smtp.SendEmailResponse{Success: true}, nil
}
