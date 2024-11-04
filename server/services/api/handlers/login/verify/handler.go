package loginVerify

import (
	"github.com/GrzegorzManiak/GOWL/pkg/crypto"
	"github.com/GrzegorzManiak/GOWL/pkg/owl"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	models2 "github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/session"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"time"
)

func handler(data *Input, ctx *gin.Context, logger *log.Logger, databaseConnection *gorm.DB) (*Output, helpers.AppError) {
	userVerify := models2.UserVerify{}
	dbErr := databaseConnection.Where("r_id = ?", data.RID).First(&userVerify)
	if dbErr.Error != nil {
		return nil, helpers.GenericError{
			Message: "User Verify not found",
			ErrCode: 400,
		}
	}

	if userVerify.CreatedAt.Unix() < time.Now().Add(-time.Second*config.Auth.MaxVerifyTimeWindow).Unix() {
		return nil, helpers.GenericError{
			Message: "User Verify expired",
			ErrCode: 400,
		}
	}

	user := models2.User{}
	dbErr = databaseConnection.Where("id = ?", userVerify.UserID).First(&user)
	if dbErr.Error != nil {
		return nil, helpers.GenericError{
			Message: "User not found",
			ErrCode: 400,
		}
	}

	owlServer, err := owl.ServerInit(user.ServerName, config.CURVE, &owl.RegistrationRequestPayload{
		U:  user.UID,
		T:  crypto.B64DecodeBytes(user.T),
		PI: crypto.B64DecodeBigInt(user.PI),
	})
	if err != nil {
		return nil, helpers.GenericError{
			Message: "User Verify, " + err.Error(),
			ErrCode: 400,
		}
	}

	clientValidate, err := parseClientValidate(data)
	if err != nil {
		return nil, helpers.GenericError{
			Message: err.Error(),
			ErrCode: 400,
		}
	}

	clientAuthInit, err := parseClientAuthInit(&userVerify, &user)
	if err != nil {
		return nil, helpers.GenericError{
			Message: err.Error(),
			ErrCode: 400,
		}
	}

	serverAuthInit, err := parseServerAuthInit(&userVerify, &user)
	if err != nil {
		return nil, helpers.GenericError{
			Message: err.Error(),
			ErrCode: 400,
		}
	}

	serverAuthValidate, err := owlServer.AuthValidate(clientAuthInit, clientValidate, serverAuthInit)
	if err != nil {
		return nil, helpers.GenericError{
			Message: err.Error(),
			ErrCode: 401,
		}
	}

	_, err = session.IssueAndSetSessionToken(ctx, user, databaseConnection)
	if err != nil {
		return nil, helpers.GenericError{
			Message: err.Error(),
			ErrCode: 400,
		}
	}

	return &Output{
		ServerKCTag: crypto.B64Encode(serverAuthValidate.Payload.ServerKCTag),

		SymmetricRootKey:     user.SymmetricRootKey,
		AsymmetricPrivateKey: user.AsymmetricPrivateKey,
		SymmetricContactsKey: user.SymmetricContactsKey,

		TotalInboundBytes:   user.TotalInboundBytes,
		TotalInboundEmails:  user.TotalInboundEmails,
		TotalOutboundBytes:  user.TotalOutboundBytes,
		TotalOutboundEmails: user.TotalOutboundEmails,
	}, nil
}
