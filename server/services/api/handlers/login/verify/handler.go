package loginVerify

import (
	"github.com/GrzegorzManiak/GOWL/pkg/crypto"
	"github.com/GrzegorzManiak/GOWL/pkg/owl"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	models2 "github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/services"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/session"
	"time"
)

func Handler(input *Input, data *services.Handler) (*Output, helpers.AppError) {
	userVerify := models2.UserVerify{}
	dbErr := data.DatabaseConnection.Where("p_id = ?", input.PID).First(&userVerify)
	if dbErr.Error != nil {
		return nil, helpers.NewUserError("Sorry! We couldn't find your account. Please try again.", "User not found")
	}

	if userVerify.CreatedAt.Unix() < time.Now().Add(-time.Second*config.Auth.MaxVerifyTimeWindow).Unix() {
		return nil, helpers.NewUserError("Sorry! Your verification window has expired. Please try again.", "Verification window expired")
	}

	user := models2.User{}
	dbErr = data.DatabaseConnection.Where("id = ?", userVerify.UserID).First(&user)
	if dbErr.Error != nil {
		return nil, helpers.NewUserError("Sorry! We couldn't find your account. Please try again.", "User not found")
	}

	owlServer, err := owl.ServerInit(user.ServerName, config.CURVE, &owl.RegistrationRequestPayload{
		U:  user.UID,
		T:  crypto.B64DecodeBytes(user.T),
		PI: crypto.B64DecodeBigInt(user.PI),
	})
	if err != nil {
		return nil, helpers.NewUserError("Sorry! We couldn't find your account. Please try again.", "User not found")
	}

	clientValidate, err := parseClientValidate(input)
	if err != nil {
		return nil, helpers.NewUserError("Sorry! We couldn't find your account. Please try again.", "User not found")
	}

	clientAuthInit, err := parseClientAuthInit(&userVerify, &user)
	if err != nil {
		return nil, helpers.NewUserError("Sorry! We couldn't find your account. Please try again.", "User not found")
	}

	serverAuthInit, err := parseServerAuthInit(&userVerify, &user)
	if err != nil {
		return nil, helpers.NewUserError("Sorry! We couldn't find your account. Please try again.", "User not found")
	}

	serverAuthValidate, err := owlServer.AuthValidate(clientAuthInit, clientValidate, serverAuthInit)
	if err != nil {
		return nil, helpers.NewUserError("Sorry! We couldn't find your account. Please try again.", "User not found")
	}

	_, err = session.IssueAndSetSessionToken(data.Context, user, data.DatabaseConnection)
	if err != nil {
		return nil, helpers.NewUserError("Sorry! We couldn't find your account. Please try again.", "User not found")
	}

	return &Output{
		ServerKCTag:   crypto.B64Encode(serverAuthValidate.Payload.ServerKCTag),
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
