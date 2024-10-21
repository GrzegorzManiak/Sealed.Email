package queue

import (
	"context"
	"gorm.io/gorm"
	"log"
	"strings"
	"sync"
	"time"
)

var rwMutex sync.RWMutex = sync.RWMutex{}

/**
I was thinking of a system that just auto-writes each request to the MySQL db
This will be our buffer.

We then have a task runner that gets x number of the next requests and executes
them.

We also need a function for refreshing a task.
*/

func PushToQueue(databaseConnection *gorm.DB, entry *Entry) (string, error) {
	rwMutex.Lock()
	defer rwMutex.Unlock()

	err := databaseConnection.Create(entry).Error
	if err != nil {
		return "", err
	}

	return entry.Uuid, nil
}

func UpdateEntry(databaseConnection *gorm.DB, entry Entry) error {
	rwMutex.Lock()
	defer rwMutex.Unlock()

	err := databaseConnection.Save(&entry).Error
	if err != nil {
		return err
	}

	return nil
}

func GetNextEntries(databaseConnection *gorm.DB, queue string, limit int) ([]Entry, error) {
	var entries []Entry
	queue = strings.ToLower(queue)
	time := time.Now().Unix()

	err := databaseConnection.
		Where("queue = ? AND status != ? AND total_attempts < permitted_attempts AND next_execution <= ?", queue, 2, time).
		Order("next_execution ASC").
		Limit(limit).
		Find(&entries).
		Error

	return entries, err
}

func MarkAsInprogress(databaseConnection *gorm.DB, uuid string) error {
	var entry Entry
	err := databaseConnection.
		Where("uuid = ?", uuid).
		First(&entry).
		Error
	if err != nil {
		return err
	}

	entry.Status = 1
	err = UpdateEntry(databaseConnection, entry)
	if err != nil {
		return err
	}

	return nil
}

func Dispatcher(
	ctx context.Context,
	databaseConnection *gorm.DB,
	queue string,
	timeout int,
	maximumWorkers int,
	worker func(entry *Entry) int8,
) {
	totalActiveWorkers := 0
	workersMutex := sync.Mutex{}
	timeoutDuration := time.Duration(timeout) * time.Second

	err := databaseConnection.AutoMigrate(&Entry{})
	if err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			if databaseConnection.Error != nil {
				log.Fatalf("Failed to get table: %v", databaseConnection.Error)
				return
			}

			if totalActiveWorkers >= maximumWorkers {
				time.Sleep(1)
				continue
			}

			entries, err := GetNextEntries(databaseConnection, queue, maximumWorkers)
			if err != nil {
				log.Printf("Failed to get entries: %v", err)
				time.Sleep(timeoutDuration)
				continue
			}

			for _, entry := range entries {
				workersMutex.Lock()
				if totalActiveWorkers >= maximumWorkers {
					workersMutex.Unlock()
					time.Sleep(timeoutDuration)
					continue
				}
				totalActiveWorkers++
				workersMutex.Unlock()

				go func(entry Entry) {
					println("Worker started", entry.Uuid)
					entry.LogAttempt()
					err := UpdateEntry(databaseConnection, entry)
					if err != nil {
						log.Printf("Failed to update entry: %v", err)
					}

					output := worker(&entry)
					entry.LogResult(output)

					err = UpdateEntry(databaseConnection, entry)
					if err != nil {
						log.Printf("Failed to update entry: %v", err)
					}

					workersMutex.Lock()
					totalActiveWorkers--
					workersMutex.Unlock()
				}(entry)
			}
		}
	}
}
