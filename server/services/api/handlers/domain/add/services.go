package domainAdd

import (
	"encoding/base64"
	"time"

	"github.com/GrzegorzManiak/NoiseBackend/config"
	models "github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/cryptography"
	"github.com/GrzegorzManiak/NoiseBackend/internal/errors"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func insertDomain(
	user *models.User,
	input *Input,
	formattedDomain string,
	databaseConnection *gorm.DB,
) (*models.UserDomain, errors.AppError) {
	kp, err := generateDKIMKeyPair()
	if err != nil {
		return &models.UserDomain{}, err
	}

	PID := helpers.GeneratePublicId(64)

	domainModel := models.UserDomain{
		PID:    PID,
		UserID: user.ID,

		Domain:   formattedDomain,
		Verified: false,
		CatchAll: false,

		DKIMKeysCreatedAt: time.Now().Unix(),
		DKIMPublicKey:     kp.EncodePublicKey(),
		DKIMPrivateKey:    kp.EncodePrivateKey(),
		TxtChallenge:      config.Domain.ChallengePrefix + "=" + helpers.GeneratePublicId(64),

		Version:             1,
		SymmetricRootKey:    input.SymmetricRootKey,
		PublicKey:           input.PublicKey,
		EncryptedPrivateKey: input.EncryptedPrivateKey,
	}

	if err := databaseConnection.Create(&domainModel); err.Error != nil {
		return &models.UserDomain{}, errors.User(
			"Domain could not be added. Please contact support if this issue persists.",
			"Failed to add domain!",
		)
	}

	return &domainModel, nil
}

func generateDKIMKeyPair() (*cryptography.RSAKeyPair, errors.AppError) {
	kp, err := cryptography.GenerateRSAKeyPair(config.Domain.DKIMKeySize)
	if err != nil {
		zap.L().Error("Error generating RSA key pair", zap.Error(err))

		return &cryptography.RSAKeyPair{}, errors.User(
			"There was an error generating the DKIM key pair. Please contact support if this issue persists.",
			"Failed to generate DKIM key pair!",
		)
	}

	return kp, nil
}

func domainAlreadyAdded(domain string, userID uint, databaseConnection *gorm.DB) bool {
	var count int64
	err := databaseConnection.
		Model(&models.UserDomain{}).
		Where("domain = ? AND user_id = ?", domain, userID).
		Count(&count)

	if err.Error != nil {
		zap.L().Error("Error checking if domain already added", zap.Error(err.Error))

		return true
	}

	return count > 0
}

func validateProofOfPossession(
	input *Input,
) bool {
	proof, err := base64.RawURLEncoding.DecodeString(input.ProofOfPossession)
	if err != nil {
		return false
	}

	decodedPublicKey, err := base64.RawURLEncoding.DecodeString(input.PublicKey)
	if err != nil {
		return false
	}

	publicKey, err := cryptography.ByteArrToECDSAPublicKey(decodedPublicKey)
	if err != nil {
		return false
	}

	if !cryptography.VerifyMessage(publicKey, input.Domain, proof) {
		return false
	}

	return true
}
