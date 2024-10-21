package queue

import (
	uuid "github.com/google/uuid"
	"gorm.io/gorm"
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
}

type EntryInterface interface {
	BeforeSave(tx *gorm.DB) (err error)
	LogAttempt(status int8)
}

// 0 - Pending, 1 - Verified, 2 - Failed, 3 - Expired
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
	return
}

func (entry *Entry) IsPending() bool {
	currentTime := time.Now().Unix()
	return (entry.Status == 0 || entry.Status == 2 || entry.NextExecution < currentTime) && entry.TotalAttempts < entry.PermittedAttempts
}

func (entry *Entry) IsExpired() bool {
	return entry.TotalAttempts >= entry.PermittedAttempts
}

func (entry *Entry) LogAttempt(status int8) {
	if status > 3 || status < 0 {
		entry.Status = 2
	} else {
		entry.Status = status
	}

	timeNow := time.Now().Unix()
	entry.LastExecution = timeNow
	entry.NextExecution = timeNow + entry.RetryInterval
	entry.TotalAttempts++
}
