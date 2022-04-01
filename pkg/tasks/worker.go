package tasks

import (
	"github.com/bluecolor/tractor/pkg/conf"
	"github.com/hibiken/asynq"
)

type Worker struct {
	*asynq.Server
}

func NewWorker(c conf.Worker) *Worker {
	worker := asynq.NewServer(
		asynq.RedisClientOpt{Addr: c.Addr},
		asynq.Config{
			Concurrency: c.Concurrency,
		},
	)
	return &Worker{
		Server: worker,
	}
}

func (w *Worker) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeSessionRun, HandleExtractionTask)
	if err := w.Run(mux); err != nil {
		return err
	}
	return nil
}
