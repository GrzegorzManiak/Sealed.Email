package list

import (
	"encoding/base64"
	"fmt"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/cryptography"
	"github.com/GrzegorzManiak/NoiseBackend/internal/errors"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strings"
)

func buildFilters(
	pagination Input,
	dbQuery *gorm.DB,
) {
	// -- Folder
	if len(pagination.Folders) > 0 {
		dbQuery = dbQuery.Where("folder IN (?)", pagination.Folders)
	}

	// -- Read
	pagination.Read = strings.ToLower(pagination.Read)
	if pagination.Read == "only" {
		dbQuery = dbQuery.Where("read = 1")
	} else if pagination.Read == "unread" {
		dbQuery = dbQuery.Where("read = 0")
	}

	// -- Sent
	pagination.Sent = strings.ToLower(pagination.Sent)
	if pagination.Sent == "in" {
		dbQuery = dbQuery.Where("sent = 1")
	} else if pagination.Sent == "out" {
		dbQuery = dbQuery.Where("sent = 0")
	}

	// -- Encrypted
	pagination.Encrypted = strings.ToLower(pagination.Encrypted)
	if pagination.Encrypted == "original" {
		dbQuery = dbQuery.Where("originally_encrypted = 1")
	} else if pagination.Encrypted == "post" {
		dbQuery = dbQuery.Where("originally_encrypted = 0")
	}

	// -- Spam
	pagination.Spam = strings.ToLower(pagination.Spam)
	if pagination.Spam == "only" {
		dbQuery = dbQuery.Where("spam = 1")
	} else if pagination.Spam == "none" {
		dbQuery = dbQuery.Where("spam = 0")
	}
}

func fetchEmails(
	user *models.User,
	pagination Input,
	databaseConnection *gorm.DB,
) (emails []*models.UserEmail, err errors.AppError) {
	emails = make([]*models.UserEmail, 0)
	dbQuery := databaseConnection.
		Table("user_emails").
		Select([]string{"read", "folder", "p_id", "domain_p_id", "`to`", "received_at", "sent", "bucket_path", "originally_encrypted", "spam"}).
		Where("user_id = ? AND domain_p_id = ?", user.ID, pagination.DomainID)

	buildFilters(pagination, dbQuery)

	if err := dbQuery.
		Limit(pagination.PerPage).
		Offset(pagination.PerPage * pagination.Page).
		Order(fmt.Sprintf("received_at %s", helpers.FormatOrderString(pagination.Order))).
		Find(&emails).Error; err != nil {
		zap.L().Debug("Failed to fetch emails", zap.Error(err))
		return nil, errors.User("The requested emails could not be found.", "Emails not found!")
	}

	return emails, nil
}

func CreateAccessKey(bucketPath string) (string, int64, error) {
	exp := helpers.GetUnixTimestamp() + 60*5 // 5 minutes
	bucketPath = fmt.Sprintf("%s:%d", bucketPath, exp)
	bytes, err := cryptography.SignMessage(&config.Session.EmailAccessPrivateKey, bucketPath)
	if err != nil {
		return "", 0, err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), exp, nil
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
