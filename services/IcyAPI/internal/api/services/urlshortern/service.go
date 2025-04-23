package urlshortern

import (
	postgresql "github.com/itsjaylen/IcyAPI/internal/api/repositories/PostgreSQL"
	redis "github.com/itsjaylen/IcyAPI/internal/api/repositories/Redis"
)

// Controller handles user-related endpoints for administrative tasks.
type Controller struct {
	Client         *redis.Client
	PostgresClient *postgresql.PostgresClient
}

// NewAdminController initializes AdminController with dependencies.
func NewURLShorternController(Client *redis.Client, postgresClient *postgresql.PostgresClient) *Controller {
	return &Controller{
		Client:         Client,
		PostgresClient: postgresClient,
	}
}