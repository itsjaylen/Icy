package workers

import (
	"IcyAPI/internal/appinit"
	"IcyAPI/internal/workers/tasks/health"
	"context"
	"fmt"
	logger "itsjaylen/IcyLogger"

	"github.com/hibiken/asynq"
)

type Worker struct { 
	App *appinit.App
}

func NewWorker(app *appinit.App) *Worker {
	return &Worker{App: app}
}

func (w *Worker) DBHealthCheckTask(ctx context.Context, task *asynq.Task) error {
	switch task.Type() {
	case "clickhouse_health_check":
		if err := health.CheckClickHouseConnection(w.App.ClickHouseClient); err != nil {
			health.UpdateHealthStatus("ClickHouse", "failed")
			logger.Error.Printf("ClickHouse Health Check Failed: %v and trying to reconnect...", err)
			w.App.ClickHouseClient.Reconnect()
			return err
		}
		health.UpdateHealthStatus("ClickHouse", "passed")

	case "postgres_health_check":
		if err := health.CheckDBConnection(w.App.PostgresClient); err != nil {
			health.UpdateHealthStatus("Postgres", "failed")
			logger.Error.Printf("Postgres Health Check Failed: %v and trying to reconnect...", err)
			w.App.PostgresClient.Reconnect()
			return err
		}
		health.UpdateHealthStatus("Postgres", "passed")

	case "minio_health_check":
		if err := health.CheckMinioConnection(w.App.MinioClient); err != nil {
			health.UpdateHealthStatus("MinIO", "failed")
			logger.Error.Printf("MinIO Health Check Failed: %v and trying to reconnect...", err)
			w.App.MinioClient.Reconnect()
			return err
		}
		health.UpdateHealthStatus("MinIO", "passed")

	case "rabbitmq_health_check":
		if err := health.CheckRabbitMQConnection(w.App.RabbitMQ); err != nil {
			health.UpdateHealthStatus("RabbitMQ", "failed")
			logger.Error.Printf("RabbitMQ Health Check Failed: %v and trying to reconnect...", err)
			w.App.RabbitMQ.Reconnect()
			return err
		}
		health.UpdateHealthStatus("RabbitMQ", "passed")

	default:
		return fmt.Errorf("unknown task type: %s", task.Type())
	}

	return nil
}
