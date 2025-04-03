// Package admin provides handlers for administrative tasks related to the API.
package admin

import (
	"net/http"

	postgresql "github.com/itsjaylen/IcyAPI/internal/api/repositories/PostgreSQL"
	redis "github.com/itsjaylen/IcyAPI/internal/api/repositories/Redis"
	utils "github.com/itsjaylen/IcyAPI/internal/utils"
	logger "itsjaylen/IcyLogger"
)

// Controller handles user-related endpoints for administrative tasks.
type Controller struct {
	Client         *redis.Client
	PostgresClient *postgresql.PostgresClient
}

// NewAdminController initializes AdminController with dependencies.
func NewAdminController(Client *redis.Client, postgresClient *postgresql.PostgresClient) *Controller {
	return &Controller{
		Client:         Client,
		PostgresClient: postgresClient,
	}
}

// HandleUserRequest example handler using dependencies.
func (c *Controller) HandleUserRequest(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context() // Use request context

	// Use Redis Set with required arguments
	err := c.Client.Set(ctx, "my_key", "my_value", 0)
	if err != nil {
		http.Error(writer, "Failed to set key in Redis", http.StatusInternalServerError)
	}

	latency, err := c.Client.Latency(ctx)
	if err != nil {
		logger.Error.Printf("Error measuring Redis latency: %v", err)
	} else {
		logger.Info.Printf("Redis latency: %v", latency)
	}

	utils.WriteJSONResponse(writer, http.StatusOK, map[string]interface{}{
		"latency": latency.Nanoseconds(),
	})
}
