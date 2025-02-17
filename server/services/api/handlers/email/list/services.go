package list

import (
	"fmt"
	"github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func fetchEmails(
	user *models.User,
	pagination Input,
	databaseConnection *gorm.DB,
) (emails []*models.UserEmail, count int64, err helpers.AppError) {
	emails = make([]*models.UserEmail, 0)
	dbQuery := databaseConnection.
		Select(helpers.BuildColumnList([]string{
			"read",
			"folder",
			"p_id",
			"domain_p_id",
			"`to`",
			"received_at",
		})).
		Where("user_id = ? AND domain_p_id = ?", user.ID, pagination.DomainID)

	if len(pagination.Folders) > 0 {
		dbQuery = dbQuery.Where("folder IN (?)", pagination.Folders)
	}

	if pagination.Read == "only" {
		dbQuery = dbQuery.Where("read = 1")
	} else if pagination.Read == "unread" {
		dbQuery = dbQuery.Where("read = 0")
	}

	if err := dbQuery.Count(&count).Error; err != nil {
		zap.L().Debug("Failed to count emails", zap.Error(err))
		return nil, 0, helpers.NewServerError("Failed to count emails.", "Email count error")
	}

	if err := dbQuery.
		Limit(pagination.PerPage).
		Offset(pagination.PerPage * pagination.Page).
		Order(fmt.Sprintf("received_at %s", helpers.FormatOrderString(pagination.Order))).
		Find(&emails).Error; err != nil {
		zap.L().Debug("Failed to fetch emails", zap.Error(err))
		return nil, 0, helpers.NewServerError("The requested emails could not be found.", "Emails not found!")
	}

	return emails, count, nil
}

func parseEmailList(
	emails []*models.UserEmail,
) *[]Email {
	emailList := make([]Email, 0)
	for _, email := range emails {
		emailList = append(emailList, Email{
			EmailID:    email.PID,
			ReceivedAt: email.ReceivedAt,
			Read:       email.Read,
			To:         email.To,
			Folder:     email.Folder,
			Spam:       email.Spam,
		})
	}
	return &emailList
}
