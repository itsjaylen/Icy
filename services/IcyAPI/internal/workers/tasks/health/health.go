package health

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	clickhouse "github.com/itsjaylen/IcyAPI/internal/api/repositories/ClickHouse"
	minobucket "github.com/itsjaylen/IcyAPI/internal/api/repositories/MinoBucket"
	postgresql "github.com/itsjaylen/IcyAPI/internal/api/repositories/PostgreSQL"
	rabbitmq "github.com/itsjaylen/IcyAPI/internal/api/repositories/RabbitMQ"
	"github.com/itsjaylen/IcyAPI/internal/utils"
)

// Status structure to hold the health check status.
type HealthStatus struct {
	Service string `json:"service"`
	Status  string `json:"status"`
}

// Global status map and a mutex for synchronization.
var (
	healthStatus      = make(map[string]string)
	healthStatusMutex sync.Mutex
)

// Update the health status with a mutex lock.
func UpdateHealthStatus(service, status string) {
	healthStatusMutex.Lock()
	defer healthStatusMutex.Unlock()
	healthStatus[service] = status
}

// ClickHouse health check.
func CheckClickHouseConnection(client *clickhouse.Client) error {
	sqlDB, err := client.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying SQL DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("ClickHouse ping failed: %w", err)
	}

	return nil
}

// PostgreSQL health check.
func CheckDBConnection(client *postgresql.PostgresClient) error {
	sqlDB, err := client.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying SQL DB: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}

// MinIO health check.
func CheckMinioConnection(client *minobucket.MinioClient) error {
	_, err := client.Client.ListBuckets(context.Background())
	if err != nil {
		return fmt.Errorf("MinIO connection failed: %w", err)
	}

	return nil
}

// RabbitMQ health check.
func CheckRabbitMQConnection(client *rabbitmq.Client) error {
	if client.Conn.IsClosed() {
		return fmt.Errorf("RabbitMQ connection is closed")
	}

	ch, err := client.Conn.Channel()
	if err != nil {
		return fmt.Errorf("RabbitMQ connection failed: %w", err)
	}
	defer ch.Close()

	return nil
}

// Healthz handler to serve status as JSON.
func HealthzHandler(writer http.ResponseWriter, _ *http.Request) {
	healthStatusMutex.Lock()
	defer healthStatusMutex.Unlock()

	status := []HealthStatus{}
	for service, statusVal := range healthStatus {
		status = append(status, HealthStatus{
			Service: service,
			Status:  statusVal,
		})
	}

	if len(status) == 0 {
		utils.WriteJSONResponse(writer, http.StatusInternalServerError, map[string]string{"error": "No health checks performed"})

		return
	}

	utils.WriteJSONResponse(writer, http.StatusOK, status)
}
