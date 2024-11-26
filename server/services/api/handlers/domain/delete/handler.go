package domainDelete

import (
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func handler(data *Input, ctx *gin.Context, userID uint, databaseConnection *gorm.DB) (*Output, helpers.AppError) {
	err := deleteDomain(userID, data.DomainID, databaseConnection)
	if err != nil {
		return nil, err
	}

	return &Output{}, nil
}
