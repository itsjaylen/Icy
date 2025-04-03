package workers

import (
	"time"

	"github.com/hibiken/asynq"
	logger "itsjaylen/IcyLogger"
)

// SchedulePeriodicHealthChecks enqueues health check tasks periodically.
func SchedulePeriodicHealthChecks(client *asynq.Client) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		// Enqueue health check tasks periodically
		tasks := []string{
			clickHouseHealthCheckTask,
			postgresHealthCheckTask,
			minioHealthCheckTask,
			rabbitMQHealthCheckTask,
		}

		for _, taskType := range tasks {
			t := asynq.NewTask(taskType, nil)
			if _, err := client.Enqueue(t, asynq.Timeout(3*time.Minute), asynq.ProcessIn(10*time.Second)); err != nil {
				logger.Error.Printf("could not enqueue %s task: %v", taskType, err)
			}
		}
	}
}
