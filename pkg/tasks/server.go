package tasks

import (
	"github.com/bluecolor/tractor/pkg/conf"
	"github.com/hibiken/asynq"
)

func NewServer(c conf.Tasks) *asynq.Server {
	return asynq.NewServer(
		asynq.RedisClientOpt{Addr: c.Addr},
		asynq.Config{
			Concurrency: c.Concurrency,
		},
	)
}
