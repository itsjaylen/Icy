// Package appinit provides functions to initialize Minio dependencies.
package appinit

import (
	"fmt"

	minobucket "github.com/itsjaylen/IcyAPI/internal/api/repositories/MinoBucket"
	config "itsjaylen/IcyConfig"
)

// InitMinio initializes a Minio client.
func InitMinio(cfg *config.AppConfig) (*minobucket.MinioClient, error) {
	addr := fmt.Sprintf("%s:%s", cfg.Minio.Host, cfg.Minio.Port)

	return minobucket.NewMinioClient(addr, cfg.Minio.AccessKey, cfg.Minio.SecretKey, false)
}
