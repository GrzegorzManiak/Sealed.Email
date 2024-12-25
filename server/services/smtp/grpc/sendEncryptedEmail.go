package grpc

import (
	"context"
	"github.com/GrzegorzManiak/NoiseBackend/proto/smtp"
)

func (s *Server) SendEncryptedEmail(ctx context.Context, encryptedEmail *smtp.EncryptedEmail) (*smtp.SendEmailResponse, error) {
	return &smtp.SendEmailResponse{Success: true}, nil
}
