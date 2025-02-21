package loginVerify

import (
	"encoding/base64"
	"github.com/GrzegorzManiak/GOWL/pkg/owl"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/errors"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/services"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/session"
	"time"
)

func Handler(input *Input, data *services.Handler) (*Output, errors.AppError) {
	userVerify := models.UserVerify{}
	dbErr := data.DatabaseConnection.Where("p_id = ?", input.PID).First(&userVerify)
	if dbErr.Error != nil {
		return nil, errors.User("Sorry! We couldn't find your account. Please try again.", "User not found")
	}

	if userVerify.CreatedAt.Unix() < time.Now().Add(-time.Second*config.Auth.MaxVerifyTimeWindow).Unix() {
		return nil, errors.User("Sorry! Your verification window has expired. Please try again.", "Verification window expired")
	}

	user := models.User{}
	dbErr = data.DatabaseConnection.Where("id = ?", userVerify.UserID).First(&user)
	if dbErr.Error != nil {
		return nil, errors.User("Sorry! We couldn't find your account. Please try again.", "User not found")
	}

	owlServer, err := owl.ServerInit(user.ServerName, config.CURVE, &owl.RegistrationRequestPayload{
		U:  user.UID,
		T:  helpers.DecodeUrlSafeBase64ToBytes(user.T),
		PI: helpers.DecodeUrlSafeBase64ToBigInt(user.PI),
	})
	if err != nil {
		return nil, errors.User("Sorry! We couldn't find your account. Please try again.", "User not found")
	}

	clientValidate, err := parseClientValidate(input)
	if err != nil {
		return nil, errors.User("Sorry! We couldn't find your account. Please try again.", "User not found")
	}

	clientAuthInit, err := parseClientAuthInit(&userVerify, &user)
	if err != nil {
		return nil, errors.User("Sorry! We couldn't find your account. Please try again.", "User not found")
	}

	serverAuthInit, err := parseServerAuthInit(&userVerify, &user)
	if err != nil {
		return nil, errors.User("Sorry! We couldn't find your account. Please try again.", "User not found")
	}

	serverAuthValidate, err := owlServer.AuthValidate(clientAuthInit, clientValidate, serverAuthInit)
	if err != nil {
		return nil, errors.User("Sorry! We couldn't find your account. Please try again.", "User not found")
	}

	_, err = session.IssueAndSetSessionToken(data.Context, user, data.DatabaseConnection)
	if err != nil {
		return nil, errors.User("Sorry! We couldn't find your account. Please try again.", "User not found")
	}

	return &Output{
		ServerKCTag:   base64.RawURLEncoding.EncodeToString(serverAuthValidate.Payload.ServerKCTag.Bytes()),
		IntegrityHash: user.IntegrityHash,

		SymmetricRootKey:     user.SymmetricRootKey,
		AsymmetricPrivateKey: user.AsymmetricPrivateKey,
		SymmetricContactsKey: user.SymmetricContactsKey,

		TotalInboundBytes:   user.TotalInboundBytes,
		TotalInboundEmails:  user.TotalInboundEmails,
		TotalOutboundBytes:  user.TotalOutboundBytes,
		TotalOutboundEmails: user.TotalOutboundEmails,
	}, nil
}
