package workers

import (
	"context"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/itsjaylen/IcyAPI/internal/appinit"
	"github.com/itsjaylen/IcyAPI/internal/workers/tasks/health"
	logger "itsjaylen/IcyLogger"
)

// Worker struct to hold dependencies.
type Worker struct {
	App *appinit.App
}

// NewWorker creates a new worker.
func NewWorker(app *appinit.App) *Worker {
	return &Worker{App: app}
}

// DBHealthCheckTask handles database health checks.
func (worker *Worker) DBHealthCheckTask(_ context.Context, task *asynq.Task) error {
	switch task.Type() {
	case "clickhouse_health_check":
		if err := health.CheckClickHouseConnection(worker.App.ClickHouseClient); err != nil {
			health.UpdateHealthStatus("ClickHouse", "failed")
			logger.Error.Printf("ClickHouse Health Check Failed: %v and trying to reconnect...", err)
			worker.App.ClickHouseClient.Reconnect()

			return err
		}
		health.UpdateHealthStatus("ClickHouse", "passed")

	case "postgres_health_check":
		if err := health.CheckDBConnection(worker.App.PostgresClient); err != nil {
			health.UpdateHealthStatus("Postgres", "failed")
			logger.Error.Printf("Postgres Health Check Failed: %v and trying to reconnect...", err)
			worker.App.PostgresClient.Reconnect()

			return err
		}
		health.UpdateHealthStatus("Postgres", "passed")

	case "minio_health_check":
		if err := health.CheckMinioConnection(worker.App.MinioClient); err != nil {
			health.UpdateHealthStatus("MinIO", "failed")
			logger.Error.Printf("MinIO Health Check Failed: %v and trying to reconnect...", err)
			worker.App.MinioClient.Reconnect()

			return err
		}
		health.UpdateHealthStatus("MinIO", "passed")

	case "rabbitmq_health_check":
		if err := health.CheckRabbitMQConnection(worker.App.RabbitMQ); err != nil {
			health.UpdateHealthStatus("RabbitMQ", "failed")
			logger.Error.Printf("RabbitMQ Health Check Failed: %v and trying to reconnect...", err)
			worker.App.RabbitMQ.Reconnect()

			return err
		}
		health.UpdateHealthStatus("RabbitMQ", "passed")

	default:
		return fmt.Errorf("unknown task type: %s", task.Type())
	}

	return nil
}
