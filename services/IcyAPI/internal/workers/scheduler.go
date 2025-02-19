package workers

import (
	logger "itsjaylen/IcyLogger"
	"time"

	"github.com/hibiken/asynq"
)

func SchedulePeriodicHealthChecks(client *asynq.Client) {
	// Set the ticker to run every 5 minutes (for example)
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
