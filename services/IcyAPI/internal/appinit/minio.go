// Package appinit provides functions to initialize Minio dependencies.
package appinit

import (
	"fmt"

	minobucket "github.com/itsjaylen/IcyAPI/internal/api/repositories/MinoBucket"
	config "itsjaylen/IcyConfig"
)

// InitMinio initializes a Minio client and ensures the required bucket exists.
func InitMinio(cfg *config.AppConfig) (*minobucket.MinioClient, error) {
	addr := fmt.Sprintf("%s:%s", cfg.Minio.Host, cfg.Minio.Port)

	client, err := minobucket.NewMinioClient(addr, cfg.Minio.AccessKey, cfg.Minio.SecretKey, false)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Minio: %w", err)
	}

	err = client.CreateBucketIfNotExists("pastebin")
	if err != nil {
		return nil, fmt.Errorf("failed to ensure bucket exists: %w", err)
	}

	return client, nil
}
