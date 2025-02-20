package admin

import (
	postgresql "IcyAPI/internal/api/repositories/PostgreSQL"
	redis "IcyAPI/internal/api/repositories/Redis"
	"IcyAPI/internal/utils"
	"context"
	"net/http"

	logger "itsjaylen/IcyLogger"
)

// AdminController handles user-related endpoints
type AdminController struct {
	RedisClient    *redis.RedisClient
	PostgresClient *postgresql.PostgresClient
}

// NewAdminController initializes AdminController with dependencies
func NewAdminController(redisClient *redis.RedisClient, postgresClient *postgresql.PostgresClient) *AdminController {
	return &AdminController{
		RedisClient:    redisClient,
		PostgresClient: postgresClient,
	}
}

// HandleUserRequest example handler using dependencies
func (c *AdminController) HandleUserRequest(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Use Redis Set with required arguments
	err := c.RedisClient.Set(ctx, "my_key", "my_value", 0)
	if err != nil {
		http.Error(w, "Failed to set key in Redis", http.StatusInternalServerError)
		return
	}

	latency, err := c.RedisClient.Latency(ctx)
	if err != nil {
		logger.Error.Printf("Error measuring Redis latency: %v", err)
	} else {
		logger.Info.Printf("Redis latency: %v", latency)
	}

	if err != nil {
		logger.Error.Printf("Error measuring Redis latency: %v", err)
	} else {
		logger.Info.Printf("Redis latency: %v", latency)
	}
	utils.WriteJSONResponse(w, http.StatusOK, map[string]interface{}{
		"latency": latency.Nanoseconds(),
	})

}
