package email

import (
	"context"
	"strings"

	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/internal/errors"
	"github.com/GrzegorzManiak/NoiseBackend/internal/service"
	smtpService "github.com/GrzegorzManiak/NoiseBackend/proto/smtp"
	"go.uber.org/zap"
)

func Send(ctx context.Context, connPool *service.Pools, email *smtpService.Email) errors.AppError {
	pool, err := connPool.GetPool(config.Etcd.SMTP.Prefix)
	if err != nil {
		zap.L().Debug("Failed to get smtp pool", zap.Error(err))

		return errors.Server("Sorry! We are unable to process your request at the moment. Please try again later.", "Failed to get smtp pool!")
	}

	conn, err := pool.GetConnection()
	if err != nil {
		zap.L().Debug("Failed to get smtp client", zap.Error(err))

		return errors.User("Sorry! We are unable to process your request at the moment. Please try again later.", "Failed to get smtp client!")
	}

	stub := smtpService.NewSmtpServiceClient(conn.Conn)
	sent, err := stub.SendEmail(ctx, email)
	zap.L().Debug("Email sent to SMTP service", zap.Any("grpc IP", conn.Conn.Target()), zap.Any("email", email))

	if err != nil {
		zap.L().Debug("Failed to queue email", zap.Error(err))

		conn.Working = false

		return errors.User(err.Error(), "Failed to queue email!")
	}

	if !sent.GetSuccess() {
		zap.L().Debug("Failed to queue email", zap.Error(err))

		conn.Working = false

		return errors.User("Sorry! We are unable to process your request at the moment. Please try again later.", "Failed to queue email!")
	}

	return nil
}

func FoldEmailBody(body string) string {
	const maxLineLength = 76 // -- Standard for email body content, wont change.
	const crlf = "\r\n"

	var result strings.Builder

	// -- split the body into lines first (by existing line breaks)
	lines := strings.Split(body, crlf)
	if len(lines) == 1 {
		// -- if there were no CRLF line breaks, try with just \n
		lines = strings.Split(body, "\n")
	}

	for _, line := range lines {
		// -- skip empty lines
		if line == "" {
			result.WriteString(crlf)
			continue
		}

		// -- while the line is longer than maxLineLength
		for len(line) > maxLineLength {

			// -- find the last space before maxLineLength
			lastSpaceIndex := strings.LastIndex(line[:maxLineLength], " ")

			if lastSpaceIndex > 0 {
				// -- write until the space and start a new line
				result.WriteString(line[:lastSpaceIndex])
				result.WriteString(crlf)
				line = line[lastSpaceIndex+1:] // +1 to skip the space

			} else {
				// -- if no space found, just cut at maxLineLength
				result.WriteString(line[:maxLineLength])
				result.WriteString(crlf)
				line = line[maxLineLength:]
			}
		}

		// -- write the remainder of the line
		if len(line) > 0 {
			result.WriteString(line)
			result.WriteString(crlf)
		}
	}

	return result.String()
}
