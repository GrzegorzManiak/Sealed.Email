package loginInit

import (
	"github.com/GrzegorzManiak/GOWL/pkg/crypto"
	"github.com/GrzegorzManiak/GOWL/pkg/owl"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/session"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func handler(data *Input, ctx *gin.Context, databaseConnection *gorm.DB) (*Output, helpers.AppError) {
	fetchedUser := models.User{}
	dbErr := databaseConnection.Where("uid = ?", data.User).First(&fetchedUser)
	if dbErr.Error != nil {
		return nil, helpers.NewUserError("Sorry! We couldn't find your account. Please try again.", "User not found")
	}

	owlServer, err := owl.ServerInit(fetchedUser.ServerName, config.CURVE, &owl.RegistrationRequestPayload{
		U:  data.User,
		T:  crypto.B64DecodeBytes(fetchedUser.T),
		PI: crypto.B64DecodeBigInt(fetchedUser.PI),
	})
	if err != nil {
		return nil, helpers.NewUserError("Sorry! We couldn't find your account. Please try again.", "User not found")
	}

	clientAuthInit, prepError := prepareClientAuthInit(data)
	if prepError != nil {
		return nil, prepError
	}

	registeredUser, parseError := parseRegisteredUser(&fetchedUser)
	if parseError != nil {
		return nil, parseError
	}

	serverAuthInit, err := owlServer.AuthInit(registeredUser, clientAuthInit)
	if err != nil {
		return nil, helpers.NewUserError("Sorry! We couldn't find your account. Please try again.", "User not found")
	}

	PID, verifyError := insertVerifyData(&fetchedUser, serverAuthInit, clientAuthInit, databaseConnection)
	if verifyError != nil {
		return nil, verifyError
	}

	session.ClearCTXSession(ctx)
	return encodeOutput(PID, serverAuthInit), nil
}
