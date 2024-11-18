package domainAdd

import (
	"fmt"
	"github.com/GrzegorzManiak/GOWL/pkg/crypto"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	models "github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/cryptography"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"gorm.io/gorm"
	"strings"
	"time"
)

func insertDomain(
	user *models.User,
	domain string,
	privateKey string,
	databaseConnection *gorm.DB,
) (models.UserDomain, helpers.AppError) {
	kp, err := generateDKIMKeyPair()
	if err != nil {
		return models.UserDomain{}, err
	}

	RID := crypto.B64Encode(crypto.GenerateKey(config.CURVE.Params().N))
	domainModel := models.UserDomain{
		RID:    RID,
		UserID: user.ID,

		Domain:   domain,
		Verified: false,

		CatchAll: false,

		DKIMKeysCreatedAt: time.Now().Unix(),
		DKIMPublicKey:     kp.EncodePublicKey(),
		DKIMPrivateKey:    kp.EncodePrivateKey(),

		Version:          1,
		EncryptedRootKey: privateKey,
	}

	if err := databaseConnection.Create(&domainModel); err.Error != nil {
		return models.UserDomain{}, helpers.GenericError{
			Message: fmt.Sprintf("Error creating domain: %v", err.Error),
			ErrCode: 500,
		}
	}

	return domainModel, nil
}

func trimDomain(domain string) (string, helpers.AppError) {
	domain = strings.Trim(domain, " .")
	if domain == "" {
		return "", helpers.GenericError{
			Message: "Domain is empty",
			ErrCode: 400,
		}
	}
	return domain, nil
}

func generateDKIMKeyPair() (*cryptography.RSAKeyPair, helpers.AppError) {
	kp, err := cryptography.GenerateRSAKeyPair(config.Auth.DKIMKeySize)
	if err != nil {
		helpers.GetLogger().Printf("Error generating RSA key pair: %v", err)
		return &cryptography.RSAKeyPair{}, helpers.GenericError{
			Message: "Error generating RSA key pair",
			ErrCode: 500,
		}
	}
	return kp, nil
}
