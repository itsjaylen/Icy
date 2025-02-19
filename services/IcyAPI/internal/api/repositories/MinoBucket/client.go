package minobucket

import (
	"IcyAPI/internal/utils"
	"context"
	"time"

	logger "itsjaylen/IcyLogger"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioClient wraps the Minio instance.
type MinioClient struct {
	Client   *minio.Client
	Endpoint string
	AccessKey string
	SecretKey string
	UseSSL   bool
}

// NewMinioClient initializes and returns a Minio client with retry logic.
func NewMinioClient(endpoint, accessKey, secretKey string, useSSL bool) (*MinioClient, error) {
	client := &MinioClient{
		Endpoint:  endpoint,
		AccessKey: accessKey,
		SecretKey: secretKey,
		UseSSL:    useSSL,
	}

	err := utils.Retry(5, 2*time.Second, client.connect)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// connect establishes a connection to Minio and performs a health check.
func (m *MinioClient) connect() error {
	client, err := minio.New(m.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(m.AccessKey, m.SecretKey, ""),
		Secure: m.UseSSL,
	})
	if err != nil {
		return err
	}

	// Perform a health check by listing buckets
	ctx := context.Background()
	_, err = client.ListBuckets(ctx)
	if err != nil {
		return err
	}

	logger.Info.Println("Connected to Minio successfully")
	m.Client = client
	return nil
}

// Reconnect attempts to reconnect to Minio using the retry utility.
func (m *MinioClient) Reconnect() {
	err := utils.Retry(5, 2*time.Second, m.connect)
	if err != nil {
		logger.Error.Println("Failed to reconnect to Minio after multiple attempts")
	} else {
		logger.Info.Println("Reconnected to Minio successfully")
	}
}
