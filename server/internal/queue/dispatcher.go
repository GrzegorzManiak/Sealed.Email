package queue

import (
	"context"
	"gorm.io/gorm"
	"log"
	"time"
)

func Dispatcher(
	ctx context.Context,
	databaseConnection *gorm.DB,
	queue *Queue,
	worker func(entry *Entry) int8,
) {
	// -- Migrate the table
	err := databaseConnection.AutoMigrate(&Entry{})
	if err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}

	go refreshPool(ctx, databaseConnection, queue, queue.BatchTimeout)
	runWorkers(ctx, queue, worker)

	ctx.Done()
}

func refreshPool(
	ctx context.Context,
	databaseConnection *gorm.DB,
	queue *Queue,
	timeout int,
) {
	for {
		select {
		// -- Context is done
		case <-ctx.Done():
			return

		//	-- Default
		default:
			// -- Check if the database connection is still alive
			if databaseConnection.Error != nil {
				log.Fatalf("Failed to get table: %v", databaseConnection.Error)
				return
			}

			// -- Get the next entries
			err := queue.GetBatch()
			if err != nil {
				log.Printf("Failed to get entries: %v", err)
				time.Sleep(time.Duration(timeout) * time.Second)
				continue
			}

			// -- Flush the queue
			err = queue.FlushQueue()
			if err != nil {
				log.Printf("Failed to flush queue: %v", err)
				time.Sleep(time.Duration(timeout) * time.Second)
				continue
			}

			// -- Batch update
			err = queue.BatchUpdate()
			if err != nil {
				log.Printf("Failed to batch update: %v", err)
				time.Sleep(time.Duration(timeout) * time.Second)
				continue
			}

			// -- Sleep
			time.Sleep(2 * time.Second)
		}
	}
}

func runWorkers(
	ctx context.Context,
	queue *Queue,
	worker func(entry *Entry) int8,
) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			job := queue.RequestWork()
			if job == nil {
				time.Sleep(time.Millisecond * 100)
				continue
			}

			go workerWrapper(job, queue, worker)
		}
	}
}

func workerWrapper(entry *Entry, queue *Queue, worker func(entry *Entry) int8) {
	println("Worker started", entry.Uuid)
	output := worker(entry)
	entry.LogStatus(output)
	queue.UpdateEntry(entry)
	queue.FinishWork(entry)
}
