package loginInit

import (
	"github.com/GrzegorzManiak/GOWL/pkg/owl"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/errors"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/services"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/session"
)

func Handler(input *Input, data *services.Handler) (*Output, errors.AppError) {
	fetchedUser := models.User{}

	dbErr := data.DatabaseConnection.Where("uid = ?", input.User).First(&fetchedUser)
	if dbErr.Error != nil {
		return nil, errors.User("Sorry! We couldn't find your account. Please try again.", "User not found")
	}

	owlServer, err := owl.ServerInit(fetchedUser.ServerName, config.CURVE, &owl.RegistrationRequestPayload{
		U:  input.User,
		T:  helpers.DecodeUrlSafeBase64ToBytes(fetchedUser.T),
		PI: helpers.DecodeUrlSafeBase64ToBigInt(fetchedUser.PI),
	})
	if err != nil {
		return nil, errors.User("Sorry! We couldn't find your account. Please try again.", "User not found")
	}

	clientAuthInit, prepError := prepareClientAuthInit(input)
	if prepError != nil {
		return nil, prepError
	}

	registeredUser, parseError := parseRegisteredUser(&fetchedUser)
	if parseError != nil {
		return nil, parseError
	}

	serverAuthInit, err := owlServer.AuthInit(registeredUser, clientAuthInit)
	if err != nil {
		return nil, errors.User("Sorry! We couldn't find your account. Please try again.", "User not found")
	}

	PID, verifyError := insertVerifyData(&fetchedUser, serverAuthInit, clientAuthInit, data.DatabaseConnection)
	if verifyError != nil {
		return nil, verifyError
	}

	session.ClearCTXSession(data.Context)

	return encodeOutput(PID, serverAuthInit), nil
}
