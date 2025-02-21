package list

import "github.com/GrzegorzManiak/NoiseBackend/services/api/session"

type Input struct {
	DomainID  string   `form:"domainID" validate:"required,PublicID"`
	Page      int      `form:"page" validate:"gte=0"`
	PerPage   int      `form:"perPage" validate:"required,gte=1,lte=30"`
	Order     string   `form:"order" validate:"omitempty,oneof=asc desc"`
	Read      string   `form:"read" validate:"omitempty,oneof=only unread all"`
	Folders   []string `form:"folders" validate:"omitempty,dive,gte=0,lte=100"`
	Sent      string   `form:"sent" validate:"omitempty,oneof=in out all"`
	Spam      string   `form:"spam" validate:"omitempty,oneof=only none all"`
	Encrypted string   `form:"encrypted" validate:"omitempty,oneof=all original post"`
}

type Output struct {
	Emails []Email `json:"emails" validate:"dive"`
	Total  int     `json:"total" validate:"gte=0"`
}

type Email struct {
	EmailID    string `json:"emailID"`
	ReceivedAt int64  `json:"receivedAt"`
	BucketPath string `json:"bucketPath"`
	Read       bool   `json:"read"`
	Folder     string `json:"folder"`
	To         string `json:"to"`
	From       string `json:"from"`
	Spam       bool   `json:"spam"`
	Sent       bool   `json:"sent"`
	AccessKey  string `json:"accessKey"`
	Expiration int64  `json:"expiration"`
}

var SessionFilter = &session.APIConfiguration{
	Allow:           []string{"default"},
	Block:           []string{},
	SessionRequired: true,
}
