package queue

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm"
)

type Queue struct {
	Name         string
	BatchTimeout int
	MaxBatchSize int

	entriesLock sync.RWMutex
	readyLock   sync.Mutex
	queueLock   sync.Mutex
	workLock    sync.Mutex

	entries *[]Entry
	ready   *[]Entry
	queue   *[]Entry
	workers int

	database *gorm.DB
}

func (q *Queue) GetBatch() (error error) {
	q.entriesLock.Lock()
	defer q.entriesLock.Unlock()

	numToFetch := q.MaxBatchSize - len(*q.entries)
	if numToFetch <= 0 {
		return nil
	}

	// -- Limit the number of entries fetched
	if float64(len(*q.queue)) > (float64(q.MaxBatchSize) * 1.5) {
		return nil
	}

	if q.database == nil {
		return errors.New("database is not initialized")
	}

	var entries []Entry
	err := q.database.
		Where("queue = ? AND status != ? AND status != ? AND total_attempts <= permitted_attempts AND next_execution <= ?", q.Name, 3, 1, time.Now().Unix()).
		Order("next_execution ASC").
		Limit(numToFetch).
		Find(&entries).
		Error
	if err != nil {
		return fmt.Errorf("failed to fetch entries: %w", err)
	}

	for i := range entries {
		entries[i].LogAttempt()
	}

	if len(entries) > 0 {
		err = q.database.Save(entries).Error
		if err != nil {
			return fmt.Errorf("failed to update entries: %w", err)
		}
	}

	*q.entries = append(*q.entries, entries...)

	return nil
}

func (q *Queue) BatchUpdate() (error error) {
	q.readyLock.Lock()
	defer q.readyLock.Unlock()

	if len(*q.ready) == 0 {
		return nil
	}

	err := q.database.Save(*q.ready).Error
	if err != nil {
		return fmt.Errorf("failed to update entries: %w", err)
	}

	*q.ready = nil

	return nil
}

func (q *Queue) FlushQueue() (error error) {
	q.queueLock.Lock()
	defer q.queueLock.Unlock()

	if len(*q.queue) == 0 {
		return nil
	}

	err := q.database.Create(*q.queue).Error
	if err != nil {
		return fmt.Errorf("failed to flush queue: %w", err)
	}

	*q.queue = nil

	return nil
}

func (q *Queue) AddEntry(entry *Entry) {
	q.queueLock.Lock()
	defer q.queueLock.Unlock()
	*q.queue = append(*q.queue, *entry)
}

func (q *Queue) UpdateEntry(entry *Entry) {
	q.readyLock.Lock()
	defer q.readyLock.Unlock()
	*q.ready = append(*q.ready, *entry)
}

func (q *Queue) RequestWork() (entry *Entry) {
	q.workLock.Lock()
	defer q.workLock.Unlock()

	if q.workers >= q.MaxBatchSize {
		return nil
	}

	if len(*q.entries) == 0 {
		return nil
	}

	entry = &(*q.entries)[0]
	*q.entries = (*q.entries)[1:]
	q.workers++

	return entry
}

func (q *Queue) FinishWork(entry *Entry) {
	q.workLock.Lock()
	defer q.workLock.Unlock()

	q.workers--
}

func NewQueue(
	databaseConnection *gorm.DB,
	queueName string,
	timeout int,
	maximumWorkers int,
) *Queue {
	queueName = strings.ToLower(queueName)

	return &Queue{
		Name:         queueName,
		BatchTimeout: timeout,
		MaxBatchSize: maximumWorkers,

		entriesLock: sync.RWMutex{},
		readyLock:   sync.Mutex{},
		queueLock:   sync.Mutex{},

		entries: &[]Entry{},
		ready:   &[]Entry{},
		queue:   &[]Entry{},

		database: databaseConnection,
	}
}
