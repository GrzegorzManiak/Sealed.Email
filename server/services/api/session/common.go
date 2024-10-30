package session

import (
	"github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"gorm.io/gorm"
	"time"
)

var (
	CurVer uint8 = 1
)

type TokenGroup struct {
	Group      string
	Ver        uint8
	DefaultEAT int64
	DefaultRAT int64
}

var (
	DefaultTokenGroup = TokenGroup{
		Group:      "default",
		Ver:        CurVer,
		DefaultEAT: 60 * 60 * 24 * 7,
	}
	SecureTokenGroup = TokenGroup{
		Group:      "secure",
		Ver:        CurVer,
		DefaultEAT: 60 * 30,
	}
)

type Header struct {
	Ver uint8 `json:"ver"`

	// -- Unless reissued, expires at is the max lifetime of the token
	//	  Refresh at is when the token has to be refreshed (Not reissued)
	ExpiresAt int64 `json:"exp"`
	CreatedAt int64 `json:"iat"`
	RefreshAt int64 `json:"rat"`

	// -- Token group is used to determine the token's purpose
	//	  Default is for general use, secure is for sensitive operations
	Group string `json:"grp"`
}

type Content struct {
	SessionID string `json:"tid"`
	UserID    string `json:"uid"`
}

type Claims struct {
	Header    Header
	Content   Content
	Signature []byte
	Session   models.Session
	Token     string
}

type APIConfiguration struct {
	SessionRequired bool
	Allow           []string
	Block           []string

	RateLimit  bool
	BucketSize int
	RefillRate float64
}

func FindSession(sessionID string, databaseConnection *gorm.DB) (models.Session, error) {
	session := models.Session{}
	err := databaseConnection.Where("sid = ?", sessionID).First(&session)
	if err.Error != nil {
		return models.Session{}, err.Error
	}
	return session, nil
}

func (claims *Claims) IsExpired() bool {
	return claims.Header.ExpiresAt < time.Now().Unix()
}

func (claims *Claims) NeedsRefresh() bool {
	return claims.Header.RefreshAt < time.Now().Unix()
}

func (claims *Claims) Filter(filter *APIConfiguration) bool {
	if len(filter.Allow) > 0 {
		for _, group := range filter.Allow {
			if claims.Header.Group == group {
				return true
			}
		}
		return false
	}

	if len(filter.Block) > 0 {
		for _, group := range filter.Block {
			if claims.Header.Group == group {
				return false
			}
		}
		return true
	}

	return true
}
