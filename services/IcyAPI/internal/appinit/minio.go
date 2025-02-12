package appinit

import (
	minobucket "IcyAPI/internal/api/repositories/MinoBucket"
	"fmt"
	config "itsjaylen/IcyConfig"
)

func InitMinio(cfg *config.AppConfig) (*minobucket.MinioClient, error) {
	addr := fmt.Sprintf("%s:%s", cfg.Minio.Host, cfg.Minio.Port)
	return minobucket.NewMinioClient(addr, cfg.Minio.AccessKey, cfg.Minio.SecretKey, false)
}
