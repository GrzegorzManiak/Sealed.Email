package server

import (
	"crypto/tls"
	"github.com/GrzegorzManiak/NoiseBackend/internal/queue"
	"gorm.io/gorm"
)

func Worker(certs *tls.Config, entry *queue.Entry, queueDatabaseConnection *gorm.DB) int8 {
	return 0
}
