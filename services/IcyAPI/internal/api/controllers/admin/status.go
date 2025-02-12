package admin

import (
	postgresql "IcyAPI/internal/api/repositories/PostgreSQL"
	redis "IcyAPI/internal/api/repositories/Redis"
	"IcyAPI/internal/utils"
	"context"
	"net/http"
	"time"
)

// UserController handles user-related endpoints
type UserController struct {
	RedisClient    *redis.RedisClient
	PostgresClient *postgresql.PostgresClient
}

// NewUserController initializes UserController with dependencies
func NewUserController(redisClient *redis.RedisClient, postgresClient *postgresql.PostgresClient) *UserController {
	return &UserController{
		RedisClient:    redisClient,
		PostgresClient: postgresClient,
	}
}

// HandleUserRequest example handler using dependencies
func (c *UserController) HandleUserRequest(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background() // Create a new context
	expiration := time.Hour     // Set expiration duration

	// Use Redis Set with required arguments
	err := c.RedisClient.Set(ctx, "key", "value", expiration)
	if err != nil {
		http.Error(w, "Failed to set key in Redis", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("User request handled"))
}

// GetStatusHandler handles the /admin/status route
func GetStatusHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"status": "ok"}
	utils.WriteJSONResponse(w, http.StatusOK, response)
}
