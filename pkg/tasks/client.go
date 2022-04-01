package tasks

import (
	"github.com/bluecolor/tractor/pkg/conf"
	"github.com/hibiken/asynq"
)

type Client struct {
	*asynq.Client
}

func NewClient(c conf.Worker) *Client {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: c.Addr})
	return &Client{
		Client: client,
	}
}
