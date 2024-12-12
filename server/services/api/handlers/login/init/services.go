package loginInit

import (
	"github.com/GrzegorzManiak/GOWL/pkg/crypto"
	"github.com/GrzegorzManiak/GOWL/pkg/owl"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	models2 "github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"gorm.io/gorm"
)

func prepareClientAuthInit(data *Input) (*owl.ClientAuthInitRequestPayload, helpers.AppError) {
	clientAuthInit := &owl.ClientAuthInitRequestPayload{
		U:  data.User,
		X1: crypto.B64DecodeBytes(data.X1),
		X2: crypto.B64DecodeBytes(data.X2),
		PI1: &crypto.SchnorrZKP{
			V: crypto.B64DecodeBytes(data.PI1_V),
			R: crypto.B64DecodeBigInt(data.PI1_R),
		},
		PI2: &crypto.SchnorrZKP{
			V: crypto.B64DecodeBytes(data.PI2_V),
			R: crypto.B64DecodeBigInt(data.PI2_R),
		},
	}

	if clientAuthInit.X1 == nil || clientAuthInit.X2 == nil || clientAuthInit.PI1 == nil || clientAuthInit.PI2 == nil {
		return nil, helpers.NewUserError("Your request is missing some data. Please try again.", "Missing data")
	}

	return clientAuthInit, nil
}

func parseRegisteredUser(fetchedUser *models2.User) (*owl.RegistrationResponse, helpers.AppError) {
	if fetchedUser == nil {
		return nil, helpers.NewUserError("Sorry! We couldn't find your account. Please try again.", "User not found")
	}

	if fetchedUser.X3 == "" || fetchedUser.PI3_V == "" || fetchedUser.PI3_R == "" {
		return nil, helpers.NewUserError("Your request is missing some data. Please try again.", "Missing data")
	}

	return &owl.RegistrationResponse{
		Payload: &owl.RegistrationResponsePayload{
			X3: crypto.B64DecodeBytes(fetchedUser.X3),
			PI3: &crypto.SchnorrZKP{
				V: crypto.B64DecodeBytes(fetchedUser.PI3_V),
				R: crypto.B64DecodeBigInt(fetchedUser.PI3_R),
			},
		},
	}, nil
}

func insertVerifyData(
	fetchedUser *models2.User,
	serverAuthInit *owl.ServerAuthInitResponse,
	clientAuthInit *owl.ClientAuthInitRequestPayload,
	databaseConnection *gorm.DB,
) (string, helpers.AppError) {
	PID := crypto.B64Encode(crypto.GenerateKey(config.CURVE.Params().N))
	newUserVerify := databaseConnection.Create(&models2.UserVerify{
		PID:      PID,
		UserID:   fetchedUser.ID,
		XPub4:    crypto.B64Encode(serverAuthInit.Xx4),
		XPri4:    crypto.B64Encode(serverAuthInit.Payload.X4),
		Beta:     crypto.B64Encode(serverAuthInit.Payload.Beta),
		PIBeta_R: crypto.B64Encode(serverAuthInit.Payload.PIBeta.R),
		PIBeta_V: crypto.B64Encode(serverAuthInit.Payload.PIBeta.V),
		PI4_R:    crypto.B64Encode(serverAuthInit.Payload.PI4.R),
		PI4_V:    crypto.B64Encode(serverAuthInit.Payload.PI4.V),
		PI1_R:    crypto.B64Encode(clientAuthInit.PI1.R),
		PI1_V:    crypto.B64Encode(clientAuthInit.PI1.V),
		PI2_R:    crypto.B64Encode(clientAuthInit.PI2.R),
		PI2_V:    crypto.B64Encode(clientAuthInit.PI2.V),
		X1:       crypto.B64Encode(clientAuthInit.X1),
		X2:       crypto.B64Encode(clientAuthInit.X2),
	})

	if newUserVerify.Error != nil {
		return "", helpers.NewServerError("Failed to create verification data. Please try again.", "Failed to create verification data")
	}

	return PID, nil
}

func encodeOutput(PID string, serverAuthInit *owl.ServerAuthInitResponse) *Output {
	return &Output{
		PID:      PID,
		X3:       crypto.B64Encode(serverAuthInit.Payload.X3),
		X4:       crypto.B64Encode(serverAuthInit.Payload.X4),
		Beta:     crypto.B64Encode(serverAuthInit.Payload.Beta),
		PI3_R:    crypto.B64Encode(serverAuthInit.Payload.PI3.R),
		PI3_V:    crypto.B64Encode(serverAuthInit.Payload.PI3.V),
		PI4_R:    crypto.B64Encode(serverAuthInit.Payload.PI4.R),
		PI4_V:    crypto.B64Encode(serverAuthInit.Payload.PI4.V),
		PIBeta_R: crypto.B64Encode(serverAuthInit.Payload.PIBeta.R),
		PIBeta_V: crypto.B64Encode(serverAuthInit.Payload.PIBeta.V),
	}
}
