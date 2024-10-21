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

func PushToQueue(databaseConnection *gorm.DB, queue string, entry *Entry) (string, error) {
	rwMutex.Lock()
	defer rwMutex.Unlock()

	queue = strings.ToLower(queue)
	err := databaseConnection.Create(entry).Error
	if err != nil {
		return "", err
	}

	return entry.Uuid, nil
}

func GetNextEntries(databaseConnection *gorm.DB, queue string, limit int) ([]EntryInterface, error) {
	rwMutex.Lock()
	defer rwMutex.Unlock()

	var entries []EntryInterface
	queue = strings.ToLower(queue)
	err := databaseConnection.Where("queue = ?", queue).Limit(limit).Find(&entries).Error
	if err != nil {
		return nil, err
	}

	return entries, nil
}

func UpdateEntry(databaseConnection *gorm.DB, entry EntryInterface) error {
	rwMutex.Lock()
	defer rwMutex.Unlock()

	err := databaseConnection.Save(entry).Error
	if err != nil {
		return err
	}

	return nil
}

func Dispatcher[T EntryInterface](ctx context.Context, databaseConnection *gorm.DB, queue string, timeout int, maximumWorkers int, worker func(entry T) int8) {
	totalActiveWorkers := 0
	workersMutex := sync.Mutex{}
	timeoutDuration := time.Duration(timeout) * time.Second

	for {
		select {
		case <-ctx.Done():
			return
		default:
			entries, err := GetNextEntries(databaseConnection, queue, maximumWorkers)
			if err != nil {
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

				go func(entry T) {
					output := worker(entry)
					entry.LogAttempt(output)
					err := UpdateEntry(databaseConnection, entry)
					if err != nil {
						log.Printf("Failed to update entry: %v", err)
					}

					workersMutex.Lock()
					totalActiveWorkers--
					workersMutex.Unlock()
				}(entry.(T))
			}

			time.Sleep(timeoutDuration)
		}
	}
}
