package session

import (
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/GrzegorzManiak/GOWL/pkg/crypto"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	models2 "github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/cryptography"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

func (claims *Claims) Sign(key *ecdsa.PrivateKey) error {

	marshaledHeader, err := json.Marshal(claims.Header)
	if err != nil {
		return err
	}

	marshaledContent, err := json.Marshal(claims.Content)
	if err != nil {
		return err
	}

	content := base64.RawURLEncoding.EncodeToString(marshaledContent)
	header := base64.RawURLEncoding.EncodeToString(marshaledHeader)
	message := fmt.Sprintf("%s.%s", header, content)

	signature, err := cryptography.SignMessage(key, message)
	if err != nil {
		return err
	}

	claims.Signature = signature
	claims.Token = fmt.Sprintf("%s.%s.%s", header, content, base64.RawURLEncoding.EncodeToString(signature))
	return nil
}

func CreateSession(user models2.User, group TokenGroup, databaseConnection *gorm.DB) (models2.Session, error) {
	unix := time.Now().Unix()
	SessionID := base64.RawURLEncoding.EncodeToString(crypto.GenerateKey(config.CURVE.Params().N).Bytes())

	session := models2.Session{
		UserID:    user.ID,
		SessionID: SessionID,
		ExpiresAt: unix + group.DefaultEAT,
		CreatedAt: unix,
		RefreshAt: unix + group.DefaultRAT,
		Group:     group.Group,
	}

	err := databaseConnection.Create(&session)
	if err.Error != nil {
		return models2.Session{}, err.Error
	}

	return session, nil
}

func CreateSessionToken(group TokenGroup, content Content, session models2.Session) (Claims, error) {
	unix := time.Now().Unix()
	header := Header{
		Ver:       group.Ver,
		ExpiresAt: unix + group.DefaultEAT,
		CreatedAt: unix,
		RefreshAt: unix + group.DefaultRAT,
		Group:     group.Group,
	}

	claims := Claims{
		Header:  header,
		Content: content,
		session: session,
	}

	zap.L().Debug("CreateSessionToken", zap.Any("claims", claims))
	err := claims.Sign(&config.Session.PrivateKey)
	if err != nil {
		return Claims{}, err
	}

	return claims, nil
}

func IssueSessionToken(user models2.User, group TokenGroup, databaseConnection *gorm.DB) (Claims, error) {
	session, err := CreateSession(user, group, databaseConnection)
	if err != nil {
		return Claims{}, err
	}

	content := Content{
		SessionID: session.SessionID,
		UserID:    user.ID,
	}

	claims, err := CreateSessionToken(group, content, session)
	if err != nil {
		return Claims{}, err
	}

	return claims, nil
}

func SetSessionCookie(ctx *gin.Context, claims Claims) {
	ctx.SetCookie(
		config.Session.CookieName,
		claims.Token,
		config.Session.CookieMaxAge,
		config.Session.CookiePath,
		config.Session.CookieDomain,
		config.Session.CookieSecure,
		true,
	)
}

func IssueAndSetSessionToken(ctx *gin.Context, user models2.User, databaseConnection *gorm.DB) (Claims, error) {
	claims, err := IssueSessionToken(user, DefaultTokenGroup, databaseConnection)
	if err != nil {
		return Claims{}, err
	}

	SetSessionCookie(ctx, claims)
	return claims, nil
}

func ClearCTXSession(ctx *gin.Context) {
	ctx.SetCookie(
		config.Session.CookieName,
		"",
		-1,
		config.Session.CookiePath,
		config.Session.CookieDomain,
		config.Session.CookieSecure,
		true,
	)
}
