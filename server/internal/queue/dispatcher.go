package queue

import (
	"context"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Dispatcher(
	ctx context.Context,
	databaseConnection *gorm.DB,
	queue *Queue,
	worker func(entry *Entry) WorkerResponse,
) {
	// -- Migrate the table
	err := databaseConnection.AutoMigrate(&Entry{})
	if err != nil {
		zap.L().Panic("failed to migrate table", zap.Error(err))
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
				zap.L().Panic("failed to get table", zap.Error(databaseConnection.Error))
			}

			// -- Get the next entries
			err := queue.GetBatch()
			if err != nil {
				zap.L().Warn("failed to get batch", zap.Error(err))
				time.Sleep(time.Duration(timeout) * time.Second)

				continue
			}

			// -- Flush the queue
			err = queue.FlushQueue()
			if err != nil {
				zap.L().Warn("failed to flush queue", zap.Error(err))
				time.Sleep(time.Duration(timeout) * time.Second)

				continue
			}

			// -- Batch update
			err = queue.BatchUpdate()
			if err != nil {
				zap.L().Warn("failed to batch update", zap.Error(err))
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
	worker func(entry *Entry) WorkerResponse,
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

func workerWrapper(entry *Entry, queue *Queue, worker func(entry *Entry) WorkerResponse) {
	output := worker(entry)
	entry.LogStatus(output)
	queue.UpdateEntry(entry)
	queue.FinishWork(entry)
}
