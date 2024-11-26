package domainList

import (
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func handler(data *Input, ctx *gin.Context, userID uint, databaseConnection *gorm.DB) (*Output, helpers.AppError) {
	domains, err := fetchDomainsByUserID(userID, *data.Pagination, databaseConnection)
	if err != nil {
		return nil, err
	}

	return &Output{
		Domains: *parseDomainList(domains),
	}, nil
}
