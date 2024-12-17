package inboxList

import (
	"github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func fetchInboxesByUserID(
	userID uint,
	domainID string,
	pagination Pagination,
	databaseConnection *gorm.DB,
) ([]*models.UserInbox, int64, helpers.AppError) {

	var count int64
	inboxes := make([]*models.UserInbox, 0)
	dbQuery := databaseConnection.
		Table("user_inboxes").
		Select("user_inboxes.id, "+
			"user_inboxes.domain_id, "+
			"user_inboxes.created_at, "+
			"user_inboxes.version, "+
			"user_inboxes.email_hash, "+
			"main.user_inboxes.encrypted_email_name, "+
			"user_inboxes.p_id").
		Joins("JOIN user_domains ON user_inboxes.domain_id = user_domains.id").
		Where("user_inboxes.user_id = ? AND user_domains.p_id = ? AND user_domains.user_id = ?", userID, domainID, userID).
		Limit(pagination.PerPage).
		Offset(pagination.PerPage * pagination.Page).
		Find(&inboxes).
		Count(&count)

	if dbQuery.Error != nil {
		return nil, 0, helpers.NewServerError(
			"The requested domains could not be found.",
			"Domains not found!",
		)
	}

	zap.L().Debug("fetchInboxesByUserID",
		zap.Any("domains", inboxes),
		zap.Any("pagination", pagination),
		zap.Any("userID", userID),
		zap.Any("domainID", domainID))

	return inboxes, count, nil
}

func parseInboxList(
	inboxes []*models.UserInbox,
) *[]Inbox {
	inboxList := make([]Inbox, 0)
	for _, inbox := range inboxes {
		inboxList = append(inboxList, Inbox{
			InboxID:            inbox.PID,
			DateAdded:          inbox.CreatedAt.Unix(),
			Version:            inbox.Version,
			EncryptedEmailName: inbox.EncryptedEmailName,
			EmailHash:          inbox.EmailHash,
		})
	}
	return &inboxList
}
