package list

import (
	"github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/errors"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/handlers/folder/create"
	"github.com/GrzegorzManiak/NoiseBackend/services/api/services"
)

func fetchFolders(data *services.Handler, domainId uint) ([]models.UserFolder, errors.AppError) {
	var folders []models.UserFolder
	dbErr := data.DatabaseConnection.
		Where("user_id = ? and user_domain_id = ?", data.User.ID, domainId).
		Find(&folders)

	if dbErr.Error != nil {
		return nil, errors.User("Sorry! We couldn't find your account. Please try again.", "User not found")
	}

	return folders, nil
}

func parseFolderList(folders []models.UserFolder) *Output {
	var result []create.Folder
	for _, folder := range folders {
		result = append(result, create.Folder{
			FolderID:            folder.PID,
			EncryptedFolderName: folder.EncryptedName,
		})
	}

	return &Output{
		Folders: result,
		Total:   len(result),
	}
}
