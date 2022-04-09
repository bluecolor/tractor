package tasks

import (
	"context"

	"net/rpc"

	"github.com/rs/zerolog/log"

	"github.com/bluecolor/tractor/pkg/conf"
	"github.com/hibiken/asynq"
)

type Worker struct {
	config conf.Worker
	*asynq.Server
}

func NewWorker(c conf.Worker) *Worker {
	worker := asynq.NewServer(
		asynq.RedisClientOpt{Addr: c.BackendAddr},
		asynq.Config{
			Concurrency: c.Concurrency,
		},
	)
	return &Worker{
		config: c,
		Server: worker,
	}
}

func (w *Worker) getFeedClient() (*rpc.Client, error) {
	return rpc.Dial("tcp", w.config.FeedBackendAddr)
}

func (w *Worker) feedClientMiddleware(h asynq.Handler) asynq.Handler {
	client, err := w.getFeedClient()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get feed backend client")
	}
	return asynq.HandlerFunc(func(ctx context.Context, t *asynq.Task) error {
		ctx = context.WithValue(ctx, "feed.client", client)
		if err := h.ProcessTask(ctx, t); err != nil {
			return err
		}
		return nil
	})
}

func (w *Worker) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeSessionRun, HandleExtractionTask)
	mux.Use(w.feedClientMiddleware)
	if err := w.Run(mux); err != nil {
		return err
	}
	return nil
}
