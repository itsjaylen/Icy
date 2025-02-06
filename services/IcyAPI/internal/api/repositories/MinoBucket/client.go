package minobucket

import (
	"context"
	logger "itsjaylen/IcyLogger"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioClient wraps the Minio instance.
type MinioClient struct {
	Client *minio.Client
}

// NewMinioClient initializes and returns a Minio client.
func NewMinioClient(endpoint, accessKey, secretKey string, useSSL bool) (*MinioClient, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	// Perform a health check by listing buckets
	ctx := context.Background()
	_, err = client.ListBuckets(ctx)
	if err != nil {
		return nil, err
	}

	logger.Info.Println("Connected to Minio successfully")
	return &MinioClient{Client: client}, nil
}
