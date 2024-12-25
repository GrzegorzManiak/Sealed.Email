package grpc

import (
	"context"
	"github.com/GrzegorzManiak/NoiseBackend/proto/smtp"
)

func (s *Server) SendPublicEmail(ctx context.Context, publicEmail *smtp.PublicEmail) (*smtp.SendEmailResponse, error) {
	return &smtp.SendEmailResponse{Success: true}, nil
}
