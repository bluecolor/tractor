package redis

import (
	"github.com/bluecolor/tractor/pkg/backends"
	"github.com/bluecolor/tractor/pkg/lib/msg"
	"github.com/bluecolor/tractor/pkg/utils"
	"github.com/go-redis/redis/v7"
)

type Config struct {
	Addr string `json:"addr"`
}

type Redis struct {
	config Config
	conn   *redis.Client
}

func New(config map[string]interface{}) (*Redis, error) {
	c := Config{}
	if err := utils.MapToStruct(config, c); err != nil {
		return nil, err
	}
	return &Redis{
		config: c,
		conn: redis.NewClient(&redis.Options{
			Addr: c.Addr,
		}),
	}, nil
}

func (r *Redis) Store(sessionID string, feedback *msg.Feedback) error {
	return nil
}

func init() {
	backends.Add("redis", func(config map[string]interface{}) (msg.FeedbackBackend, error) {
		return New(config)
	})
}
