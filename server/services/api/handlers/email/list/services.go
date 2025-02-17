package list

import (
	"fmt"
	"github.com/GrzegorzManiak/GOWL/pkg/crypto"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/cryptography"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strings"
)

func fetchEmails(
	user *models.User,
	pagination Input,
	databaseConnection *gorm.DB,
) (emails []*models.UserEmail, err helpers.AppError) {
	emails = make([]*models.UserEmail, 0)
	dbQuery := databaseConnection.
		Table("user_emails").
		Select([]string{"read", "folder", "p_id", "domain_p_id", "`to`", "received_at", "sent", "bucket_path"}).
		Where("user_id = ? AND domain_p_id = ?", user.ID, pagination.DomainID)

	if len(pagination.Folders) > 0 {
		dbQuery = dbQuery.Where("folder IN (?)", pagination.Folders)
	}

	pagination.Read = strings.ToLower(pagination.Read)
	if pagination.Read == "only" {
		dbQuery = dbQuery.Where("read = 1")
	} else if pagination.Read == "unread" {
		dbQuery = dbQuery.Where("read = 0")
	}

	pagination.Sent = strings.ToLower(pagination.Sent)
	if pagination.Sent == "in" {
		dbQuery = dbQuery.Where("sent = 1")
	} else if pagination.Sent == "out" {
		dbQuery = dbQuery.Where("sent = 0")
	}

	if err := dbQuery.
		Limit(pagination.PerPage).
		Offset(pagination.PerPage * pagination.Page).
		Order(fmt.Sprintf("received_at %s", helpers.FormatOrderString(pagination.Order))).
		Find(&emails).Error; err != nil {
		zap.L().Debug("Failed to fetch emails", zap.Error(err))
		return nil, helpers.NewServerError("The requested emails could not be found.", "Emails not found!")
	}

	return emails, nil
}

func CreateAccessKey(bucketPath string) (string, int64, error) {
	exp := helpers.GetUnixTimestamp() + 60*5 // 5 minutes
	bucketPath += fmt.Sprintf(":%d", exp)
	bytes, err := cryptography.SignMessage(&config.Session.EmailAccessPrivateKey, bucketPath)
	if err != nil {
		return "", 0, err
	}
	return crypto.B64Encode(bytes), exp, nil
}

func ParseEmail(
	email *models.UserEmail,
) *Email {
	account, exp, err := CreateAccessKey(email.BucketPath)
	if err != nil {
		zap.L().Debug("Failed to create access key", zap.Error(err))
		return nil
	}

	return &Email{
		EmailID:    email.PID,
		BucketPath: email.BucketPath,
		ReceivedAt: email.ReceivedAt,
		Read:       email.Read,
		To:         email.To,
		Folder:     email.Folder,
		Spam:       email.Spam,
		Sent:       email.Sent,
		AccessKey:  account,
		Expiration: exp,
	}
}

func parseEmailList(
	emails []*models.UserEmail,
) *Output {
	emailList := make([]Email, 0)
	count := 0

	for _, email := range emails {
		emailData := ParseEmail(email)
		if emailData == nil {
			continue
		}
		emailList = append(emailList, *emailData)
		count++
	}

	return &Output{
		Emails: emailList,
		Total:  count,
	}
}
