package domainList

import (
	"github.com/GrzegorzManiak/NoiseBackend/database/primary/models"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func fetchDomainsByUserID(
	userID uint,
	pagination Pagination,
	databaseConnection *gorm.DB,
) ([]*models.UserDomain, helpers.AppError) {
	domains := make([]*models.UserDomain, 0)
	dbQuery := databaseConnection.Where("user_id = ?", userID).Limit(pagination.PerPage).Offset(pagination.PerPage * pagination.Page).Find(&domains)
	if dbQuery.Error != nil {
		return nil, helpers.NewServerError("The requested domains could not be found.", "Domains not found!")
	}

	zap.L().Debug("fetchDomainsByUserID", zap.Any("domains", domains), zap.Any("pagination", pagination), zap.Any("userID", userID))
	return domains, nil
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
