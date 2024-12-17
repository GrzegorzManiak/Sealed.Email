package domainList

import (
	"github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func fetchDomainsByUserID(
	user *models.User,
	pagination Input,
	databaseConnection *gorm.DB,
) ([]*models.UserDomain, int64, helpers.AppError) {

	var count int64
	domains := make([]*models.UserDomain, 0)
	dbQuery := databaseConnection.
		Select("p_id,"+
			"domain,"+
			"verified,"+
			"created_at,"+
			"catch_all,"+
			"version").
		Where("user_id = ?", user.ID).
		Limit(pagination.PerPage).
		Offset(pagination.PerPage * pagination.Page).
		Find(&domains).
		Count(&count)

	if dbQuery.Error != nil {
		return nil, 0, helpers.NewServerError(
			"The requested domains could not be found.",
			"Domains not found!",
		)
	}

	zap.L().Debug("fetchDomainsByUserID",
		zap.Any("domains", domains),
		zap.Any("pagination", pagination),
		zap.Any("userID", user.ID))

	return domains, count, nil
}

func parseDomainList(
	domains []*models.UserDomain,
) *[]Domain {
	domainList := make([]Domain, 0)
	for _, domain := range domains {
		domainList = append(domainList, Domain{
			DomainID:  domain.PID,
			Domain:    domain.Domain,
			Verified:  domain.Verified,
			DateAdded: domain.CreatedAt.Unix(),
			CatchAll:  domain.CatchAll,
			Version:   domain.Version,
		})
	}
	return &domainList
}
