package register

import (
	"fmt"
	"github.com/GrzegorzManiak/GOWL/pkg/owl"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/internal/cryptography"
	"github.com/GrzegorzManiak/NoiseBackend/internal/errors"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/services"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/session"
)

func Handler(input *Input, data *services.Handler) (*Output, errors.AppError) {

	proof := helpers.DecodeUrlSafeBase64ToBytes(input.Proof)
	publicKey, err := cryptography.ByteArrToECDSAPublicKey(helpers.DecodeUrlSafeBase64ToBytes(input.PublicKey))
	if err != nil {
		return nil, errors.User(fmt.Sprintf("Error converting public key: %v", err), "Oops! Something went wrong")
	}

	if !cryptography.VerifyMessage(publicKey, input.User, proof) {
		return nil, errors.User("Uh oh! Looks like your proof is invalid. Please try again.", "Invalid key proof")
	}

	owlServer, err := owl.ServerInit("NoiseEmailServer>V1.0.0", config.CURVE, &owl.RegistrationRequestPayload{
		U:  input.User,
		T:  helpers.DecodeUrlSafeBase64ToBytes(input.T),
		PI: helpers.DecodeUrlSafeBase64ToBigInt(input.PI),
	})

	if err != nil {
		return nil, errors.User(fmt.Sprintf("Error initializing server: %v", err), "Oops! Something went wrong")
	}

	if usernameExists(input.User, data.DatabaseConnection) {
		return nil, errors.User("Sorry, that username is already taken. Please try another.", "Username taken")
	}

	registeredUser := owlServer.RegisterUser()
	_, dbErr := registerUser(input, registeredUser, data.DatabaseConnection)
	if dbErr != nil {
		return nil, dbErr
	}
	session.ClearCTXSession(data.Context)

	return &Output{
		Message: "User registered successfully",
	}, nil
}
