package email

import (
	"context"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/GrzegorzManiak/NoiseBackend/internal/service"
	smtpService "github.com/GrzegorzManiak/NoiseBackend/proto/smtp"
	"go.uber.org/zap"
)

func Email(ctx context.Context, connPool *service.Pools, email *smtpService.Email) helpers.AppError {
	pool, err := connPool.GetPool(config.Etcd.SMTP.Prefix)
	if err != nil {
		zap.L().Debug("Failed to get smtp pool", zap.Error(err))
		return helpers.NewServerError("Sorry! We are unable to process your request at the moment. Please try again later.", "Failed to get smtp pool!")
	}

	conn, err := pool.GetConnection()
	if err != nil {
		zap.L().Debug("Failed to get smtp client", zap.Error(err))
		return helpers.NewServerError("Sorry! We are unable to process your request at the moment. Please try again later.", "Failed to get smtp client!")
	}

	stub := smtpService.NewSmtpServiceClient(conn.Conn)
	sent, err := stub.SendEmail(ctx, email)
	zap.L().Debug("Email sent to SMTP service", zap.Any("grpc IP", conn.Conn.Target()), zap.Any("email", email))

	if err != nil {
		zap.L().Debug("Failed to queue email", zap.Error(err))
		conn.Working = false
		return helpers.NewServerError(err.Error(), "Failed to queue email!")
	}

	if !sent.Success {
		zap.L().Debug("Failed to queue email", zap.Error(err))
		conn.Working = false
		return helpers.NewServerError("Sorry! We are unable to process your request at the moment. Please try again later.", "Failed to queue email!")
	}

	return nil
}
