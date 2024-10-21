package session

import (
	"encoding/json"
	"fmt"
	"github.com/GrzegorzManiak/GOWL/pkg/crypto"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/internal/cryptography"
	"gorm.io/gorm"
	"strings"
)

func ParseToken(token string) (Header, Content, []byte, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return Header{}, Content{}, nil, fmt.Errorf("invalid token format")
	}

	headerBytes := crypto.B64DecodeBytes(parts[0])
	contentBytes := crypto.B64DecodeBytes(parts[1])
	signature := crypto.B64DecodeBytes(parts[2])

	header := Header{}
	content := Content{}

	err := json.Unmarshal(headerBytes, &header)
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

	content := crypto.B64Encode(contentBytes)
	header := crypto.B64Encode(headerBytes)
	message := fmt.Sprintf("%s.%s", header, content)

	return cryptography.VerifyMessage(&config.Session.PublicKey, message, claims.Signature)
}

func RefreshSessionToken(claims Claims, databaseConnection *gorm.DB) (Claims, error) {
	session, err := FindSession(claims.Content.SessionID, databaseConnection)
	if err != nil {
		return Claims{}, err
	}

	if session.SessionID == "" {
		return Claims{}, fmt.Errorf("session not found")
	}

	if session.Revoked {
		return Claims{}, fmt.Errorf("session has been revoked")
	}

	newClaims, err := CreateSessionToken(DefaultTokenGroup, claims.Content, session)
	if err != nil {
		return Claims{}, err
	}

	return newClaims, nil
}
