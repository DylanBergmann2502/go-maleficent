// internal/pkg/jobs/system.go
package jobs

import (
	"fmt"
	"log/slog"

	"github.com/DylanBergmann2502/go-maleficent/config"
	"github.com/hibiken/asynq"
)

func redisOptions(config *config.Config) asynq.RedisClientOpt {
	return asynq.RedisClientOpt{
		Addr:     fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	}
}

type Worker struct {
	client *Client
	server *asynq.Server
	mux    *asynq.ServeMux
}

func NewWorker(config *config.Config, logger *slog.Logger) (*Worker, error) {
	redis := redisOptions(config)
	client, err := NewClient(config)
	if err != nil {
		return nil, err
	}

	mux := asynq.NewServeMux()
	RegisterHandlers(mux, logger)
	server := asynq.NewServer(redis, asynq.Config{
		Concurrency: config.Jobs.Concurrency,
		Queues:      map[string]int{"default": 1},
		Logger:      asynqLogger{logger: logger},
	})

	return &Worker{client: client, server: server, mux: mux}, nil
}

func (w *Worker) Client() *Client {
	return w.client
}

func (w *Worker) Run() error {
	return w.server.Run(w.mux)
}

func (w *Worker) Close() error {
	w.server.Shutdown()
	return w.client.Close()
}

type Scheduler struct {
	scheduler *asynq.Scheduler
}

type Client struct {
	client *asynq.Client
}

func NewClient(config *config.Config) (*Client, error) {
	client := asynq.NewClient(redisOptions(config))
	if err := client.Ping(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &Client{client: client}, nil
}

func (c *Client) Enqueue(task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return c.client.Enqueue(task, opts...)
}

func (c *Client) Close() error {
	return c.client.Close()
}

func NewScheduler(config *config.Config, logger *slog.Logger) (*Scheduler, error) {
	redis := redisOptions(config)
	client, err := NewClient(config)
	if err != nil {
		return nil, err
	}
	if err := client.Close(); err != nil {
		return nil, fmt.Errorf("failed to close Redis connection: %w", err)
	}

	scheduler := asynq.NewScheduler(redis, &asynq.SchedulerOpts{
		Logger: asynqLogger{logger: logger},
	})
	if err := RegisterSchedules(scheduler, config); err != nil {
		scheduler.Shutdown()
		return nil, fmt.Errorf("failed to register example schedule: %w", err)
	}

	return &Scheduler{scheduler: scheduler}, nil
}

func (s *Scheduler) Run() error {
	return s.scheduler.Run()
}

func (s *Scheduler) Close() {
	s.scheduler.Shutdown()
}
