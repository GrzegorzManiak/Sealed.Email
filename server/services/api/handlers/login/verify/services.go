package loginVerify

import (
	"github.com/GrzegorzManiak/GOWL/pkg/crypto"
	"github.com/GrzegorzManiak/GOWL/pkg/owl"
	models2 "github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/errors"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
)

func parseClientValidate(data *Input) (*owl.ClientAuthValidateRequestPayload, errors.AppError) {
	clientValidate := owl.ClientAuthValidateRequestPayload{
		ClientKCTag: helpers.DecodeUrlSafeBase64ToBigInt(data.ClientKCTag),
		Alpha:       helpers.DecodeUrlSafeBase64ToBytes(data.Alpha),
		PIAlpha: &crypto.SchnorrZKP{
			V: helpers.DecodeUrlSafeBase64ToBytes(data.PIAlpha_V),
			R: helpers.DecodeUrlSafeBase64ToBigInt(data.PIAlpha_R),
		},
		R: helpers.DecodeUrlSafeBase64ToBigInt(data.R),
	}

	return &clientValidate, nil
}

func parseClientAuthInit(userVerify *models2.UserVerify, user *models2.User) (*owl.ClientAuthInitRequestPayload, errors.AppError) {
	clientAuthInit := owl.ClientAuthInitRequestPayload{
		U:  user.UID,
		X1: helpers.DecodeUrlSafeBase64ToBytes(userVerify.X1),
		X2: helpers.DecodeUrlSafeBase64ToBytes(userVerify.X2),
		PI1: &crypto.SchnorrZKP{
			V: helpers.DecodeUrlSafeBase64ToBytes(userVerify.PI1_V),
			R: helpers.DecodeUrlSafeBase64ToBigInt(userVerify.PI1_R),
		},
		PI2: &crypto.SchnorrZKP{
			V: helpers.DecodeUrlSafeBase64ToBytes(userVerify.PI2_V),
			R: helpers.DecodeUrlSafeBase64ToBigInt(userVerify.PI2_R),
		},
	}

	return &clientAuthInit, nil
}

func parseServerAuthInit(userVerify *models2.UserVerify, user *models2.User) (*owl.ServerAuthInitResponse, errors.AppError) {
	serverAuthInit := owl.ServerAuthInitResponse{
		Xx4: helpers.DecodeUrlSafeBase64ToBigInt(userVerify.XPub4),
		Payload: &owl.ServerAuthInitResponsePayload{
			X3: helpers.DecodeUrlSafeBase64ToBytes(user.X3),
			X4: helpers.DecodeUrlSafeBase64ToBytes(userVerify.XPri4),
			PI3: &crypto.SchnorrZKP{
				V: helpers.DecodeUrlSafeBase64ToBytes(user.PI3_V),
				R: helpers.DecodeUrlSafeBase64ToBigInt(user.PI3_R),
			},
			PI4: &crypto.SchnorrZKP{
				V: helpers.DecodeUrlSafeBase64ToBytes(userVerify.PI4_V),
				R: helpers.DecodeUrlSafeBase64ToBigInt(userVerify.PI4_R),
			},
			Beta: helpers.DecodeUrlSafeBase64ToBytes(userVerify.Beta),
			PIBeta: &crypto.SchnorrZKP{
				V: helpers.DecodeUrlSafeBase64ToBytes(userVerify.PIBeta_V),
				R: helpers.DecodeUrlSafeBase64ToBigInt(userVerify.PIBeta_R),
			},
		},
	}

	return &serverAuthInit, nil
}
