package queue

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WorkerResponse int

type Entry struct {
	gorm.Model
	Uuid  string `gorm:"unique"`
	RefID string `gorm:"unique, index"`

	LastExecution     int64
	NextExecution     int64
	RetryInterval     int64
	TotalAttempts     int64
	PermittedAttempts int64
	Status            WorkerResponse
	Queue             string
	Data              string
}

type EntryData interface {
	Marshal() (string, error)
}

const (
	Pending  WorkerResponse = iota
	Verified WorkerResponse = iota
	Failed   WorkerResponse = iota
	Expired  WorkerResponse = iota
)

// 3 - Expired.
func (entry *Entry) BeforeSave(tx *gorm.DB) (err error) {
	if entry.Status > 3 || entry.Status < 0 {
		entry.Status = 2
	}

	if entry.Uuid == "" {
		entryUuid, err := uuid.NewUUID()
		if err != nil {
			return err
		}

		entry.Uuid = entryUuid.String()
	}

	entry.Queue = strings.ToLower(entry.Queue)

	return
}

func (entry *Entry) LogAttempt() {
	timeNow := time.Now().Unix()
	entry.LastExecution = timeNow
	entry.NextExecution = timeNow + entry.RetryInterval
	entry.TotalAttempts++

	if entry.TotalAttempts >= entry.PermittedAttempts {
		entry.Status = 3
	}
}

func (entry *Entry) LogStatus(status WorkerResponse) {
	if entry.TotalAttempts >= entry.PermittedAttempts {
		entry.Status = 3
	} else if status >= 0 && status <= 3 {
		entry.Status = status
	} else {
		entry.Status = 2
	}
}

func (entry *Entry) String() string {
	return fmt.Sprintf("Entry{Uuid: %s, LastExecution: %d, NextExecution: %d, RetryInterval: %d, TotalAttempts: %d, PermittedAttempts: %d, Status: %d, Queue: %s}", entry.Uuid, entry.LastExecution, entry.NextExecution, entry.RetryInterval, entry.TotalAttempts, entry.PermittedAttempts, entry.Status, entry.Queue)
}

func Initiate(maxAttempts int64, retryInterval int64, queue string, data EntryData) (*Entry, error) {
	entry := Entry{}

	dataString, err := data.Marshal()
	if err != nil {
		return nil, err
	}

	entry.PermittedAttempts = maxAttempts
	entry.RetryInterval = retryInterval
	entry.Queue = queue
	timeNow := time.Now().Unix()
	entry.NextExecution = timeNow + (60 * 0)
	entry.Data = dataString

	return &entry, err
}
