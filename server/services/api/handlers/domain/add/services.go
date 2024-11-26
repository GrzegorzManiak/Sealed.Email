package domainAdd

import (
	"fmt"
	"github.com/GrzegorzManiak/GOWL/pkg/crypto"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	models "github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/cryptography"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

func insertDomain(
	user *models.User,
	domain string,
	privateKey string,
	databaseConnection *gorm.DB,
) (*models.UserDomain, helpers.AppError) {
	kp, err := generateDKIMKeyPair()
	if err != nil {
		return &models.UserDomain{}, err
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
		TxtChallenge:      config.Domain.ChallengePrefix + "=" + crypto.B64Encode(crypto.GenerateKey(config.CURVE.Params().N)),

		Version:          1,
		EncryptedRootKey: privateKey,
	}

	if err := databaseConnection.Create(&domainModel); err.Error != nil {
		return &models.UserDomain{}, helpers.GenericError{
			Message: fmt.Sprintf("Error creating domain: %v", err.Error),
			ErrCode: 500,
		}
	}

	return &domainModel, nil
}

func generateDKIMKeyPair() (*cryptography.RSAKeyPair, helpers.AppError) {
	kp, err := cryptography.GenerateRSAKeyPair(config.Domain.DKIMKeySize)
	if err != nil {
		zap.L().Error("Error generating RSA key pair", zap.Error(err))
		return &cryptography.RSAKeyPair{}, helpers.GenericError{
			Message: "Error generating RSA key pair",
			ErrCode: 500,
		}
	}
	return kp, nil
}

func domainAlreadyAdded(domain string, userID uint, databaseConnection *gorm.DB) bool {
	var count int64
	err := databaseConnection.Model(&models.UserDomain{}).Where("domain = ? AND user_id = ?", domain, userID).Count(&count)
	if err.Error != nil {
		zap.L().Error("Error checking if domain already added", zap.Error(err.Error))
		return true
	}
	return count > 0
}
