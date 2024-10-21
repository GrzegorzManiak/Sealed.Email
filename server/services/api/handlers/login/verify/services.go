package loginVerify

import (
	"github.com/GrzegorzManiak/GOWL/pkg/crypto"
	"github.com/GrzegorzManiak/GOWL/pkg/owl"
	models2 "github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
)

func parseClientValidate(data *Input) (*owl.ClientAuthValidateRequestPayload, helpers.AppError) {
	clientValidate := owl.ClientAuthValidateRequestPayload{
		ClientKCTag: crypto.B64DecodeBigInt(data.ClientKCTag),
		Alpha:       crypto.B64DecodeBytes(data.Alpha),
		PIAlpha: &crypto.SchnorrZKP{
			V: crypto.B64DecodeBytes(data.PIAlpha_V),
			R: crypto.B64DecodeBigInt(data.PIAlpha_R),
		},
		R: crypto.B64DecodeBigInt(data.R),
	}

	return &clientValidate, nil
}

func parseClientAuthInit(userVerify *models2.UserVerify, user *models2.User) (*owl.ClientAuthInitRequestPayload, helpers.AppError) {
	clientAuthInit := owl.ClientAuthInitRequestPayload{
		U:  user.UID,
		X1: crypto.B64DecodeBytes(userVerify.X1),
		X2: crypto.B64DecodeBytes(userVerify.X2),
		PI1: &crypto.SchnorrZKP{
			V: crypto.B64DecodeBytes(userVerify.PI1_V),
			R: crypto.B64DecodeBigInt(userVerify.PI1_R),
		},
		PI2: &crypto.SchnorrZKP{
			V: crypto.B64DecodeBytes(userVerify.PI2_V),
			R: crypto.B64DecodeBigInt(userVerify.PI2_R),
		},
	}

	return &clientAuthInit, nil
}

func parseServerAuthInit(userVerify *models2.UserVerify, user *models2.User) (*owl.ServerAuthInitResponse, helpers.AppError) {
	serverAuthInit := owl.ServerAuthInitResponse{
		Xx4: crypto.B64DecodeBigInt(userVerify.XPub4),
		Payload: &owl.ServerAuthInitResponsePayload{
			X3: crypto.B64DecodeBytes(user.X3),
			X4: crypto.B64DecodeBytes(userVerify.XPri4),
			PI3: &crypto.SchnorrZKP{
				V: crypto.B64DecodeBytes(user.PI3_V),
				R: crypto.B64DecodeBigInt(user.PI3_R),
			},
			PI4: &crypto.SchnorrZKP{
				V: crypto.B64DecodeBytes(userVerify.PI4_V),
				R: crypto.B64DecodeBigInt(userVerify.PI4_R),
			},
			Beta: crypto.B64DecodeBytes(userVerify.Beta),
			PIBeta: &crypto.SchnorrZKP{
				V: crypto.B64DecodeBytes(userVerify.PIBeta_V),
				R: crypto.B64DecodeBigInt(userVerify.PIBeta_R),
			},
		},
	}

	return &serverAuthInit, nil
}
