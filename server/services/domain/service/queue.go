package service

import (
	"fmt"
	database "github.com/GrzegorzManiak/NoiseBackend/database/domain/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/queue"
	"gorm.io/gorm"
	"sync"
	"time"
)

var primaryDatabaseMutex = sync.Mutex{}
var QueueName = "Domain Verification Queue"

func Worker(entry *queue.Entry, primaryDatabaseConnection *gorm.DB) int8 {

	data, err := database.UnmarshalVerificationQueue(entry.Data)
	if err != nil {
		return 2
	}

	time.Sleep(time.Second * 5)

	// -- Lock when we will be writing to primary database
	primaryDatabaseMutex.Lock()
	defer primaryDatabaseMutex.Unlock()

	fmt.Println("Executing Queue: ", entry.String())
	fmt.Println("Data: ", data)
	return 0
}
