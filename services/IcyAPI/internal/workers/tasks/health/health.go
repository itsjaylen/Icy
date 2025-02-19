package health

import (
	clickhouse "IcyAPI/internal/api/repositories/ClickHouse"
	minobucket "IcyAPI/internal/api/repositories/MinoBucket"
	postgresql "IcyAPI/internal/api/repositories/PostgreSQL"
	rabbitmq "IcyAPI/internal/api/repositories/RabbitMQ"
	"IcyAPI/internal/utils"
	"context"
	"fmt"
	"net/http"
	"sync"
)

// Status structure to hold the health check status
type HealthStatus struct {
	Service string `json:"service"`
	Status  string `json:"status"`
}

// Global status map and a mutex for synchronization
var healthStatus = make(map[string]string)
var healthStatusMutex sync.Mutex

// Update the health status with a mutex lock
func UpdateHealthStatus(service, status string) {
	healthStatusMutex.Lock()
	defer healthStatusMutex.Unlock()
	healthStatus[service] = status
}

// ClickHouse health check
func CheckClickHouseConnection(client *clickhouse.ClickHouseClient) error {
	sqlDB, err := client.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying SQL DB: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("ClickHouse ping failed: %v", err)
	}
	return nil
}

// PostgreSQL health check
func CheckDBConnection(client *postgresql.PostgresClient) error {
	sqlDB, err := client.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying SQL DB: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %v", err)
	}
	return nil
}

// MinIO health check
func CheckMinioConnection(client *minobucket.MinioClient) error {
	_, err := client.Client.ListBuckets(context.Background())
	if err != nil {
		return fmt.Errorf("MinIO connection failed: %v", err)
	}
	return nil
}

// RabbitMQ health check
func CheckRabbitMQConnection(client *rabbitmq.RabbitMQClient) error {
	if client.Conn.IsClosed() {
		return fmt.Errorf("RabbitMQ connection is closed")
	}

	ch, err := client.Conn.Channel()
	if err != nil {
		return fmt.Errorf("RabbitMQ connection failed: %v", err)
	}
	defer ch.Close()
	return nil
}

// Healthz handler to serve status as JSON
func HealthzHandler(w http.ResponseWriter, r *http.Request) {
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
		utils.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "No health checks performed"})
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, status)
}
