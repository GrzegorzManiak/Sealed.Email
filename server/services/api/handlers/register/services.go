package register

import (
	"github.com/GrzegorzManiak/GOWL/pkg/crypto"
	"github.com/GrzegorzManiak/GOWL/pkg/owl"
	"github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal"
	"gorm.io/gorm"
)

func registerUser(
	data *Input,
	registeredUser *owl.RegistrationResponse,
	databaseConnection *gorm.DB,
) (*models.User, internal.AppError) {
	user := &models.User{
		UID: data.User,

		ServerName: "NoiseEmailServer>V1.0.0",
		T:          data.T,
		PI:         data.PI,
		X3:         crypto.B64Encode(registeredUser.Payload.X3),
		PI3_V:      crypto.B64Encode(registeredUser.Payload.PI3.V),
		PI3_R:      crypto.B64Encode(registeredUser.Payload.PI3.R),

		SymmetricRootKey:     data.EncryptedRootKey,
		AsymmetricPublicKey:  data.PublicKey,
		AsymmetricPrivateKey: data.EncryptedPrivateKey,
	}

	dbInsert := databaseConnection.Create(user)
	if dbInsert.Error != nil {
		return nil, internal.GenericError{
			Message: dbInsert.Error.Error(),
			ErrCode: 400,
		}
	}

	return user, nil
}
