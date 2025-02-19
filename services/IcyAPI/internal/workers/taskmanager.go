package workers

import (
	"IcyAPI/internal/appinit"
	"fmt"
	logger "itsjaylen/IcyLogger"
	"time"

	"github.com/hibiken/asynq"
)

const (
	clickHouseHealthCheckTask = "clickhouse_health_check"
	postgresHealthCheckTask   = "postgres_health_check"
	minioHealthCheckTask      = "minio_health_check"
	rabbitMQHealthCheckTask   = "rabbitmq_health_check"
)

type TaskManagerController struct {
	App    *appinit.App
	Worker *Worker
}

func NewTaskManagerController(app *appinit.App) *TaskManagerController {
	return &TaskManagerController{
		App:    app,
		Worker: NewWorker(app),
	}
}

// sets up the task manager for background tasks
func SetupTaskManager(app *appinit.App) {
	logger.Info.Println("Setting up task manager...")
	mux := asynq.NewServeMux()
	taskManagerController := NewTaskManagerController(app)

	mux.HandleFunc(clickHouseHealthCheckTask, taskManagerController.Worker.DBHealthCheckTask)
	mux.HandleFunc(postgresHealthCheckTask, taskManagerController.Worker.DBHealthCheckTask)
	mux.HandleFunc(minioHealthCheckTask, taskManagerController.Worker.DBHealthCheckTask)
	mux.HandleFunc(rabbitMQHealthCheckTask, taskManagerController.Worker.DBHealthCheckTask)

	redisOpts := asynq.RedisClientOpt{
		Addr: fmt.Sprintf("%s:%s", app.Cfg.Redis.Host, app.Cfg.Redis.Port),
	}

	server := asynq.NewServer(redisOpts, asynq.Config{
		Concurrency: 5,
		Queues: map[string]int{
			"critical": 6,
			"default":  3,
			"low":      1,
		},
	})

	client := asynq.NewClient(redisOpts)

	go func() {
		if err := server.Run(mux); err != nil {
			logger.Debug.Fatalf("Error running server: %v", err)
		}
	}()

	time.Sleep(2 * time.Second)
	go SchedulePeriodicHealthChecks(client)
}
