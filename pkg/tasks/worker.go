package tasks

import (
	"context"

	"github.com/bluecolor/tractor/pkg/conf"
	"github.com/bluecolor/tractor/pkg/lib/msg"
	"github.com/bluecolor/tractor/pkg/tasks/middleware/feedbackend"
	"github.com/hibiken/asynq"
)

type Worker struct {
	*asynq.Server
	feedBackends []msg.FeedBackend
}

func NewWorker(c conf.Worker) *Worker {

	worker := asynq.NewServer(
		asynq.RedisClientOpt{Addr: c.BackendAddr},
		asynq.Config{
			Concurrency: c.Concurrency,
		},
	)
	fb, err := feedbackend.New(c.BackendAddr, c.DB)
	if err != nil {
		panic(err)
	}
	return &Worker{
		Server:       worker,
		feedBackends: []msg.FeedBackend{fb},
	}
}

func (w *Worker) loggingMiddleware(h asynq.Handler) asynq.Handler {
	return asynq.HandlerFunc(func(ctx context.Context, t *asynq.Task) error {
		ctx = context.WithValue(ctx, "feedBackends", w.feedBackends)
		err := h.ProcessTask(ctx, t)
		if err != nil {
			return err
		}
		return nil
	})
}

func (w *Worker) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeSessionRun, HandleExtractionTask)
	mux.Use(w.loggingMiddleware)
	if err := w.Run(mux); err != nil {
		return err
	}
	return nil
}
