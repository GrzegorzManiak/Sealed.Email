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
) ([]*models.UserInbox, helpers.AppError) {
	inboxes := make([]*models.UserInbox, 0)
	dbQuery := databaseConnection.
		Table("user_inboxes").
		Select("user_inboxes.id, user_inboxes.domain_id, user_inboxes.created_at, user_inboxes.version, user_inboxes.email_hash, user_inboxes.p_id").
		Joins("JOIN user_domains ON user_inboxes.domain_id = user_domains.id").
		Where("user_inboxes.user_id = ? AND user_domains.p_id = ? AND user_domains.user_id = ?", userID, domainID, userID).
		Limit(pagination.PerPage).
		Offset(pagination.PerPage * pagination.Page).
		Find(&inboxes)

	if dbQuery.Error != nil {
		return nil, helpers.NewServerError(
			"The requested domains could not be found.",
			"Domains not found!",
		)
	}

	zap.L().Debug("fetchInboxesByUserID",
		zap.Any("domains", inboxes),
		zap.Any("pagination", pagination),
		zap.Any("userID", userID),
		zap.Any("domainID", domainID))

	return inboxes, nil
}

func parseInboxList(
	inboxes []*models.UserInbox,
) *[]Inbox {
	inboxList := make([]Inbox, 0)
	for _, inbox := range inboxes {
		inboxList = append(inboxList, Inbox{
			InboxID:   inbox.PID,
			InboxName: inbox.EmailHash,
			DateAdded: inbox.CreatedAt.Unix(),
			Version:   inbox.Version,
		})
	}
	return &inboxList
}
