package email

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/GrzegorzManiak/GOWL/pkg/crypto"
	"github.com/GrzegorzManiak/NoiseBackend/internal/cryptography"
	"github.com/GrzegorzManiak/NoiseBackend/internal/validation"
	smtpService "github.com/GrzegorzManiak/NoiseBackend/proto/smtp"
)

type EncryptedInbox struct {
	DisplayName       string `json:"displayName"       validate:"lte=200"`
	EmailHash         string `json:"emailHash"         validate:"required,email"`
	PublicKey         string `json:"publicKey"         validate:"EncodedP256Key"`
	EncryptedEmailKey string `json:"encryptedEmailKey" validate:"EncodedEncryptedKey"`
}

type EncryptionKey struct {
	EmailKey  string
	PublicKey string
}

func (i EncryptionKey) String() string {
	return fmt.Sprintf("%s <%s>", i.PublicKey, i.EmailKey)
}

func (i EncryptedInbox) BasicInbox() Inbox {
	return Inbox{
		DisplayName: i.DisplayName,
		Email:       i.EmailHash,
	}
}

func ReMapEncryptedInboxes(inboxes []EncryptedInbox) []Inbox {
	var result []Inbox
	for _, inbox := range inboxes {
		result = append(result, inbox.BasicInbox())
	}

	return result
}

func ConvertToInboxKeys(inboxes ...[]EncryptedInbox) []*smtpService.InboxKeys {
	var result []*smtpService.InboxKeys

	for _, inbox := range inboxes {
		for _, encryptedInbox := range inbox {
			result = append(result, &smtpService.InboxKeys{
				DisplayName:       encryptedInbox.DisplayName,
				PublicKey:         encryptedInbox.PublicKey,
				EmailHash:         validation.NormalizeEmail(encryptedInbox.EmailHash),
				EncryptedEmailKey: encryptedInbox.EncryptedEmailKey,
			})
		}
	}

	return result
}

func StringifyInboxKeys(inboxes []EncryptedInbox) string {
	var result []string

	for _, inbox := range inboxes {
		stringified := fmt.Sprintf("<%s:%s>", inbox.PublicKey, inbox.EncryptedEmailKey)
		result = append(result, stringified)
	}

	return strings.Join(result, ",")
}

func CreateInboxKey() ([]byte, error) {
	return cryptography.NewKey(cryptography.DefaultKeyLength)
}

func EncryptEmailKey(emailKey []byte, publicKey string) (*EncryptionKey, error) {
	decodedKey, err := base64.RawURLEncoding.DecodeString(publicKey)
	if err != nil {
		return nil, err
	}

	ecdsaPublicKey, err := cryptography.ByteArrToECDSAPublicKey(decodedKey)
	if err != nil {
		return nil, err
	}

	cipherText, err := cryptography.AsymEncrypt(ecdsaPublicKey, emailKey)
	if err != nil {
		return nil, err
	}

	return &EncryptionKey{
		EmailKey:  base64.RawURLEncoding.EncodeToString(cipherText),
		PublicKey: publicKey,
	}, nil
}

func HashInboxEmail(email string) (string, error) {
	email = validation.NormalizeEmail(email)

	domain, err := validation.ExtractDomainFromEmail(email)
	if err != nil {
		return "", err
	}

	user := strings.SplitN(email, "@", 2)[0]
	email = fmt.Sprintf("%s@%s", user, domain)
	hashedEmail := crypto.Hash(email)
	encodedEmail := base64.RawURLEncoding.EncodeToString(hashedEmail.Bytes())

	return fmt.Sprintf("%s@%s", encodedEmail, domain), nil
}
