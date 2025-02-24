package session

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/internal/cryptography"
	"gorm.io/gorm"
)

func ParseToken(token string) (Header, Content, []byte, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return Header{}, Content{}, nil, errors.New("invalid token format")
	}

	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return Header{}, Content{}, nil, err
	}

	contentBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return Header{}, Content{}, nil, err
	}

	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return Header{}, Content{}, nil, err
	}

	header := Header{}
	content := Content{}

	err = json.Unmarshal(headerBytes, &header)
	if err != nil {
		return Header{}, Content{}, nil, err
	}

	err = json.Unmarshal(contentBytes, &content)
	if err != nil {
		return Header{}, Content{}, nil, err
	}

	return header, content, signature, nil
}

func ParseSessionToken(token string) (Claims, error) {
	claims := Claims{}

	header, content, signature, err := ParseToken(token)
	if err != nil {
		return Claims{}, err
	}

	claims.Header = header
	claims.Content = content
	claims.Signature = signature
	claims.Token = token

	return claims, nil
}

func (claims *Claims) Verify() bool {
	headerBytes, err := json.Marshal(claims.Header)
	if err != nil {
		return false
	}

	contentBytes, err := json.Marshal(claims.Content)
	if err != nil {
		return false
	}

	content := base64.RawURLEncoding.EncodeToString(contentBytes)
	header := base64.RawURLEncoding.EncodeToString(headerBytes)
	message := fmt.Sprintf("%s.%s", header, content)

	return cryptography.VerifyMessage(&config.Session.PublicKey, message, claims.Signature)
}

func RefreshSessionToken(claims Claims, databaseConnection *gorm.DB) (Claims, error) {
	session, err := FindSession(claims.Content.SessionID, databaseConnection)
	if err != nil {
		return Claims{}, err
	}

	if session.SessionID == "" {
		return Claims{}, errors.New("session not found")
	}

	if session.Revoked {
		return Claims{}, errors.New("session has been revoked")
	}

	newClaims, err := CreateSessionToken(DefaultTokenGroup, claims.Content, session)
	if err != nil {
		return Claims{}, err
	}

	return newClaims, nil
}
