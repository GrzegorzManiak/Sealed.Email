package service

import (
	database "github.com/GrzegorzManiak/NoiseBackend/database/domain/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/queue"
	"gorm.io/gorm"
	"math/rand"
	"sync"
	"time"
)

var primaryDatabaseMutex = sync.Mutex{}
var QueueName = "Domain Verification Queue"

func Worker(entry *queue.Entry, primaryDatabaseConnection *gorm.DB) int8 {

	_, err := database.UnmarshalVerificationQueue(entry.Data)
	if err != nil {
		return 2
	}

	// -- Sleep for 5 - 10 seconds
	time.Sleep(time.Duration(5+rand.Intn(5)) * time.Second)

	// -- Return 1/2
	random := rand.Intn(2)
	if random == 0 {
		return 1
	} else {
		return 2
	}
}
