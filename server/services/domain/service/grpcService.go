package service

import (
	"github.com/GrzegorzManiak/NoiseBackend/internal/queue"
	"github.com/GrzegorzManiak/NoiseBackend/proto/domain"
	"gorm.io/gorm"
)

type Server struct {
	domain.UnimplementedDomainServiceServer
	QueueDatabaseConnection *gorm.DB
	Queue                   *queue.Queue
}
