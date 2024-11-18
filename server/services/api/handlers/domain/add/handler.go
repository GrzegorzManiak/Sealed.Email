package domainAdd

import (
	"fmt"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/gin-gonic/gin"
	"log"
)

func handler(data *Input, ctx *gin.Context, logger *log.Logger) (*Output, helpers.AppError) {

	fmt.Println("Adding: ", data.Domain)

	return nil, nil
}
