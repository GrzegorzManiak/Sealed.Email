package register

import (
	"github.com/GrzegorzManiak/GOWL/pkg/crypto"
	"github.com/GrzegorzManiak/GOWL/pkg/owl"
	"github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"gorm.io/gorm"
)

func usernameExists(username string, databaseConnection *gorm.DB) bool {
	var user models.User
	databaseConnection.Where("uid = ?", username).First(&user)
	return user.UID != ""
}

func registerUser(
	data *Input,
	registeredUser *owl.RegistrationResponse,
	databaseConnection *gorm.DB,
) (*models.User, helpers.AppError) {
	user := &models.User{
		UID: data.User,

		ServerName: "NoiseEmailServer>V1.0.0",
		T:          data.T,
		PI:         data.PI,
		X3:         crypto.B64Encode(registeredUser.Payload.X3),
		PI3_V:      crypto.B64Encode(registeredUser.Payload.PI3.V),
		PI3_R:      crypto.B64Encode(registeredUser.Payload.PI3.R),

		IntegrityHash: data.IntegrityHash,

		SymmetricRootKey:     data.EncryptedRootKey,
		AsymmetricPublicKey:  data.PublicKey,
		AsymmetricPrivateKey: data.EncryptedPrivateKey,
		SymmetricContactsKey: data.EncryptedContactsKey,

		TotalInboundEmails:  0,
		TotalInboundBytes:   0,
		TotalOutboundEmails: 0,
		TotalOutboundBytes:  0,
	}

	dbInsert := databaseConnection.Create(user)
	if dbInsert.Error != nil {
		return nil, helpers.NewServerError("Error inserting user into database", "Oops! Something went wrong")
	}

	return user, nil
}
