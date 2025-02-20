package session

import (
	"github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/errors"
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
		DefaultEAT: 60 * 60 * 24 * 8, // -- 8 days
		DefaultRAT: 60 * 5,           // -- 5 minutes
	}
	SecureTokenGroup = TokenGroup{
		Group:      "secure",
		Ver:        CurVer,
		DefaultEAT: 60 * 30, // -- 30 minutes
		DefaultRAT: 60 * 5,  // -- 5 minutes
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
	UserID    uint   `json:"uid"`
}

type Claims struct {
	Header    Header
	Content   Content
	Signature []byte
	session   models.Session
	Token     string
}

type APIConfiguration struct {
	SessionRequired bool
	Bypass          bool
	SelfResponse    bool
	Allow           []string
	Block           []string
}

func FindSession(sessionID string, databaseConnection *gorm.DB) (models.Session, error) {
	session := models.Session{}
	err := databaseConnection.Where("session_id = ?", sessionID).First(&session)
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

func (claims *Claims) FetchUser(databaseConnection *gorm.DB) (models.User, errors.AppError) {
	user := models.User{}
	err := databaseConnection.
		Where("id = ?", claims.Content.UserID).
		First(&user)

	if err.Error != nil {
		return models.User{}, errors.User(
			"Sorry! We couldn't find your account. Please try again.",
			"User not found",
		)
	}

	return user, nil
}
