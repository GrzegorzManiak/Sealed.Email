package client

import (
	"github.com/GrzegorzManiak/NoiseBackend/internal/queue"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Worker(entry *queue.Entry, queueDatabaseConnection *gorm.DB) int8 {
	zap.L().Debug("Processing smtp queue", zap.Any("entry", entry))
	return 1
}
