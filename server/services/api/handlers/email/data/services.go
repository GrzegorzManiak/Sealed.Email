package data

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/GrzegorzManiak/NoiseBackend/config"
	"github.com/GrzegorzManiak/NoiseBackend/internal/cryptography"
	"github.com/GrzegorzManiak/NoiseBackend/internal/helpers"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"io"
)

func fetchEmailData(input *Input, minioClient *minio.Client, writer *gin.ResponseWriter) helpers.AppError {
	object, err := minioClient.GetObject(context.Background(), "emails", input.BucketPath, minio.GetObjectOptions{})
	if err != nil {
		return helpers.NewServerError("Failed to fetch email data", "Failed to fetch email data")
	}
	defer object.Close()

	if _, err := io.Copy(*writer, object); err != nil {
		return helpers.NewServerError("Failed to write email data", "Failed to write email data")
	}

	return nil
}

func validateAccessKey(input *Input) bool {
	if input.Expiration < helpers.GetUnixTimestamp() {
		return false
	}
	bucketPath := fmt.Sprintf("%s:%d", input.BucketPath, input.Expiration)
	decodedAccessKey, err := base64.RawURLEncoding.DecodeString(input.AccessKey)
	if err != nil {
		return false
	}
	return cryptography.VerifyMessage(&config.Session.EmailAccessPublicKey, bucketPath, decodedAccessKey)
}
