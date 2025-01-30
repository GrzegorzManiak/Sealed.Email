package register

import (
	"fmt"
	"github.com/GrzegorzManiak/GOWL/pkg/crypto"
	"github.com/GrzegorzManiak/GOWL/pkg/owl"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/internal/cryptography"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/services"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/session"
)

func Handler(input *Input, data *services.Handler) (*Output, helpers.AppError) {

	proof := crypto.B64DecodeBytes(input.Proof)
	publicKey, err := cryptography.ByteArrToECDSAPublicKey(config.CURVE, crypto.B64DecodeBytes(input.PublicKey))
	if err != nil {

		return nil, helpers.NewServerError(fmt.Sprintf("Error converting public key: %v", err), "Oops! Something went wrong")
	}

	if !cryptography.VerifyMessage(publicKey, input.User, proof) {
		return nil, helpers.NewUserError("Uh oh! Looks like your proof is invalid. Please try again.", "Invalid key proof")
	}

	owlServer, err := owl.ServerInit("NoiseEmailServer>V1.0.0", config.CURVE, &owl.RegistrationRequestPayload{
		U:  input.User,
		T:  crypto.B64DecodeBytes(input.T),
		PI: crypto.B64DecodeBigInt(input.PI),
	})

	if err != nil {
		return nil, helpers.NewServerError(fmt.Sprintf("Error initializing server: %v", err), "Oops! Something went wrong")
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
