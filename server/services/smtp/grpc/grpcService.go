package grpc

import (
	"github.com/GrzegorzManiak/NoiseBackend/internal/queue"
	"github.com/GrzegorzManiak/NoiseBackend/proto/smtp"
	"gorm.io/gorm"
)

type Server struct {
	smtp.UnimplementedSmtpServiceServer
	MainDatabaseConnection  *gorm.DB
	QueueDatabaseConnection *gorm.DB

	InboundQueue  *queue.Queue
	OutboundQueue *queue.Queue
}
