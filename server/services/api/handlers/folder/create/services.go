package create

import (
	"github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/errors"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/services"
	"gorm.io/gorm"
)

func CheckDomainOwnership(databaseConnection *gorm.DB, domainId string, userId uint) (bool, uint) {
	var domain models.UserDomain
	response := databaseConnection.Where("p_id = ? AND user_id = ?", domainId, userId).First(&domain)
	if response.Error != nil {
		return false, 0
	}
	return true, domain.ID
}

func createFolder(input *Input, data *services.Handler, domainId uint) (*Folder, errors.AppError) {
	folder := &models.UserFolder{
		UserId:        data.User.ID,
		PID:           helpers.GeneratePublicId(32),
		EncryptedName: input.EncryptedFolderName,
		UserDomainID:  domainId,
	}

	if err := data.DatabaseConnection.Create(&folder); err.Error != nil {
		return nil, errors.User(
			"Could not create the folder.",
			"Folder creation failed!",
		)
	}

	return &Folder{
		FolderID:            folder.PID,
		EncryptedFolderName: folder.EncryptedName,
	}, nil
}
