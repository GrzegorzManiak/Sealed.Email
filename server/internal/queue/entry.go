package queue

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Entry struct {
	gorm.Model
	Uuid              string `gorm:"unique"`
	LastExecution     int64
	NextExecution     int64
	RetryInterval     int64
	TotalAttempts     int64
	PermittedAttempts int64
	Status            int8
	Queue             string
	Data              string
}

type EntryData interface {
	Marshal() (string, error)
}

// BeforeSave
// 0 - Pending,
// 1 - Verified,
// 2 - Failed,
// 3 - Expired
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

func (entry *Entry) LogStatus(status int8) {
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
	entry.NextExecution = timeNow
	entry.Data = dataString

	return &entry, err
}
