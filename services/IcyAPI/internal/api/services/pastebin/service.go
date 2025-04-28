package pastebin

import (
	"github.com/go-playground/validator"
	minobucket "github.com/itsjaylen/IcyAPI/internal/api/repositories/MinoBucket"
	postgresql "github.com/itsjaylen/IcyAPI/internal/api/repositories/PostgreSQL"
	"github.com/microcosm-cc/bluemonday"

	config "itsjaylen/IcyConfig"
)

// Controller handles user-related endpoints for administrative tasks.
type Controller struct {
	PostgresClient *postgresql.PostgresClient
	MinoClient      *minobucket.MinioClient
	Validator *validator.Validate
	Sanitizer   *bluemonday.Policy
	Cfg              *config.AppConfig
}

func NewPasteBinController(PostgresClient *postgresql.PostgresClient, MinoClient *minobucket.MinioClient) *Controller {
	return &Controller{
		PostgresClient: PostgresClient,
		MinoClient:      MinoClient,
		Validator: validator.New(),
		Sanitizer:   bluemonday.UGCPolicy(),
		Cfg:                &config.AppConfig{},
	}
}
