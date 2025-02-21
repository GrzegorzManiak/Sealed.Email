package delete

import (
	"github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"gorm.io/gorm"
)

func deleteFolder(databaseConnection *gorm.DB, domainId, userId uint, folderPid string) error {
	response := databaseConnection.Where("p_id = ? AND user_id = ? AND user_domain_id = ?", folderPid, userId, domainId).Delete(&models.UserFolder{})
	if response.Error != nil {
		return response.Error
	}
	return nil
}
