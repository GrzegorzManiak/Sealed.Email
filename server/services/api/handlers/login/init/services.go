package loginInit

import (
	"github.com/GrzegorzManiak/GOWL/pkg/crypto"
	"github.com/GrzegorzManiak/GOWL/pkg/owl"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	models2 "github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal"
	"gorm.io/gorm"
)

func prepareClientAuthInit(data *Input) (*owl.ClientAuthInitRequestPayload, internal.AppError) {
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
		return nil, internal.GenericError{
			Message: "clientAuthInit fields are not properly initialized",
			ErrCode: 400,
		}
	}

	return clientAuthInit, nil
}

func parseRegisteredUser(fetchedUser *models2.User) (*owl.RegistrationResponse, internal.AppError) {
	if fetchedUser == nil {
		return nil, internal.GenericError{
			Message: "fetchedUser is nil",
			ErrCode: 400,
		}
	}

	if fetchedUser.X3 == "" || fetchedUser.PI3_V == "" || fetchedUser.PI3_R == "" {
		return nil, internal.GenericError{
			Message: "fetchedUser fields are not properly initialized",
			ErrCode: 400,
		}
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
) (string, internal.AppError) {
	RID := crypto.B64Encode(crypto.GenerateKey(config.CURVE.Params().N))
	newUserVerify := databaseConnection.Create(&models2.UserVerify{
		RID:      RID,
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
		return "", internal.GenericError{
			Message: newUserVerify.Error.Error(),
			ErrCode: 400,
		}
	}

	return RID, nil
}

func encodeOutput(RID string, serverAuthInit *owl.ServerAuthInitResponse) *Output {
	return &Output{
		RID:      RID,
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
