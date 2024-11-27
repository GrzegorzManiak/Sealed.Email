package devSession

import "github.com/GrzegorzManiak/NoiseBackend/services/api/session"

type Input struct {
}

type Output struct {
	Token     string `json:"token"`
	Signature []byte `json:"signature"`

	ContentUserID    uint   `json:"contentUserID"`
	ContentSessionID string `json:"contentSessionID"`

	HeaderVersion uint8  `json:"headerVersion"`
	HeaderGroup   string `json:"headerGroup"`
	HeaderEat     int64  `json:"headerExpiresAt"`
	HeaderCat     int64  `json:"headerCreatedAt"`
	HeaderRat     int64  `json:"headerRefreshAt"`
}

var SessionFilter = &session.APIConfiguration{
	Allow:           []string{"default"},
	Block:           []string{},
	SessionRequired: true,
}
