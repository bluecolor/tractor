package tasks

import (
	"github.com/bluecolor/tractor/pkg/conf"
	"github.com/hibiken/asynq"
)

type Client struct {
	config *conf.Tasks
	client *asynq.Client
}

func NewClient(c conf.Tasks) *Client {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: c.Addr})
	return &Client{
		config: &c,
		client: client,
	}
}
func (c *Client) Close() error {
	return c.client.Close()
}
func (c *Client) Client() *asynq.Client {
	return c.client
}
func (c *Client) Enqueue(task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	return c.client.Enqueue(task, opts...)
}
