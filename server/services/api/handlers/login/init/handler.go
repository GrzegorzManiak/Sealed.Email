package loginInit

import (
	"github.com/GrzegorzManiak/GOWL/pkg/crypto"
	"github.com/GrzegorzManiak/GOWL/pkg/owl"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
)

func handler(data *Input, ctx *gin.Context, logger *log.Logger, databaseConnection *gorm.DB) (*Output, helpers.AppError) {
	fetchedUser := models.User{}
	dbErr := databaseConnection.Where("uid = ?", data.User).First(&fetchedUser)
	if dbErr.Error != nil {
		return nil, helpers.GenericError{
			Message: "User not found",
			ErrCode: 400,
		}
	}

	owlServer, err := owl.ServerInit(fetchedUser.ServerName, config.CURVE, &owl.RegistrationRequestPayload{
		U:  data.User,
		T:  crypto.B64DecodeBytes(fetchedUser.T),
		PI: crypto.B64DecodeBigInt(fetchedUser.PI),
	})
	if err != nil {
		return nil, helpers.GenericError{
			Message: "User Init, " + err.Error(),
			ErrCode: 400,
		}
	}

	clientAuthInit, prepError := prepareClientAuthInit(data)
	if prepError != nil {
		return nil, prepError
	}

	registeredUser, parseError := parseRegisteredUser(&fetchedUser)
	if parseError != nil {
		return nil, parseError
	}

	serverAuthInit, err := owlServer.AuthInit(registeredUser, clientAuthInit)
	if err != nil {
		return nil, helpers.GenericError{
			Message: err.Error(),
			ErrCode: 400,
		}
	}

	RID, verifyError := insertVerifyData(&fetchedUser, serverAuthInit, clientAuthInit, databaseConnection)
	if verifyError != nil {
		return nil, verifyError
	}

	return encodeOutput(RID, serverAuthInit), nil
}
