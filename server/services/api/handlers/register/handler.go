package register

import (
	"fmt"
	"github.com/GrzegorzManiak/GOWL/pkg/crypto"
	"github.com/GrzegorzManiak/GOWL/pkg/owl"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/internal/cryptography"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/session"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func handler(data *Input, ctx *gin.Context, databaseConnection *gorm.DB) (*Output, helpers.AppError) {

	proof := crypto.B64DecodeBytes(data.Proof)
	publicKey, err := cryptography.ByteArrToECDSAPublicKey(config.CURVE, crypto.B64DecodeBytes(data.PublicKey))
	if err != nil {

		return nil, helpers.NewServerError(fmt.Sprintf("Error converting public key: %v", err), "Oops! Something went wrong")
	}

	if !cryptography.VerifyMessage(publicKey, data.User, proof) {
		return nil, helpers.NewUserError("Uh oh! Looks like your proof is invalid. Please try again.", "Invalid key proof")
	}

	owlServer, err := owl.ServerInit("NoiseEmailServer>V1.0.0", config.CURVE, &owl.RegistrationRequestPayload{
		U:  data.User,
		T:  crypto.B64DecodeBytes(data.T),
		PI: crypto.B64DecodeBigInt(data.PI),
	})

	if err != nil {
		return nil, helpers.NewServerError(fmt.Sprintf("Error initializing server: %v", err), "Oops! Something went wrong")
	}

	registeredUser := owlServer.RegisterUser()
	_, dbErr := registerUser(data, registeredUser, databaseConnection)
	if dbErr != nil {
		return nil, dbErr
	}
	session.ClearCTXSession(ctx)

	return &Output{
		Message: "User registered successfully",
	}, nil
}
