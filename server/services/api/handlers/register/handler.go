package register

import (
	"github.com/GrzegorzManiak/GOWL/pkg/crypto"
	"github.com/GrzegorzManiak/GOWL/pkg/owl"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/internal"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/session"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
)

func handler(data *Input, ctx *gin.Context, logger *log.Logger, databaseConnection *gorm.DB) (*Output, internal.AppError) {

	proof := crypto.B64DecodeBytes(data.Proof)
	publicKey, err := internal.ByteArrToECDSAPublicKey(config.CURVE, crypto.B64DecodeBytes(data.PublicKey))
	if err != nil {

		return nil, internal.GenericError{
			Message: err.Error(),
			ErrCode: 400,
		}
	}

	if !internal.VerifyMessage(publicKey, data.User, proof) {
		return nil, internal.GenericError{
			Message: "Invalid proof",
			ErrCode: 400,
		}
	}

	owlServer, err := owl.ServerInit("NoiseEmailServer>V1.0.0", config.CURVE, &owl.RegistrationRequestPayload{
		U:  data.User,
		T:  crypto.B64DecodeBytes(data.T),
		PI: crypto.B64DecodeBigInt(data.PI),
	})

	if err != nil {
		return nil, internal.GenericError{
			Message: err.Error(),
			ErrCode: 400,
		}
	}

	registeredUser := owlServer.RegisterUser()
	newUser, dbErr := registerUser(data, registeredUser, databaseConnection)
	if dbErr != nil {
		return nil, dbErr
	}

	_, err = session.IssueAndSetSessionToken(ctx, *newUser, databaseConnection)
	if err != nil {
		return nil, internal.GenericError{
			Message: err.Error(),
			ErrCode: 400,
		}
	}

	return &Output{
		Message: "User registered successfully",
	}, nil
}
