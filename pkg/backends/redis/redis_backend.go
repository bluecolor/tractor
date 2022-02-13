package redis

import (
	"strings"

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
	if err := utils.MapToStruct(config, &c); err != nil {
		return nil, err
	}
	return &Redis{
		config: c,
		conn: redis.NewClient(&redis.Options{
			Addr: c.Addr,
		}),
	}, nil
}
func (r *Redis) Close() error {
	return r.conn.Close()
}
func (r *Redis) TrySetRunning(sessionID string) error {
	session, err := r.GetSession(sessionID)
	if err != nil {
		return err
	}
	_, ok := session[StatusKey]
	if !ok {
		r.SetStatus(sessionID, "running")
	}
	return nil
}
func (r *Redis) GetSession(sessionID string) (map[string]string, error) {
	key := getSessionKey(sessionID)
	return r.conn.HGetAll(key).Result()
}
func (r *Redis) SetStatus(sessionID string, status string) error {
	key := getSessionKey(sessionID)
	return r.conn.HSet(key, StatusKey, strings.ToLower(status)).Err()
}
func (r *Redis) SetInputProgress(sessionID string, progress int) error {
	key := getSessionKey(sessionID)
	return r.conn.HSet(key, InputProgressKey, progress).Err()
}
func (r *Redis) SetOutputProgress(sessionID string, progress int) error {
	key := getSessionKey(sessionID)
	return r.conn.HSet(key, OutputProgressKey, progress).Err()
}
func (r *Redis) SetError(sessionID string, err error) error {
	key := getSessionKey(sessionID)
	return r.conn.HSet(key, ErrorKey, err.Error()).Err()
}
func (r *Redis) SetWarning(sessionID string, warning string) error {
	key := getSessionKey(sessionID)
	return r.conn.HSet(key, WarningKey, warning).Err()
}
func (r *Redis) Store(sessionID string, feedback *msg.Feedback) error {
	if err := r.TrySetRunning(sessionID); err != nil {
		return err
	}
	if err := r.Publish(sessionID, feedback); err != nil {
		return err
	}
	switch feedback.Type {
	case msg.Error:
		r.SetStatus(sessionID, msg.Error.String())
		r.SetError(sessionID, feedback.Error())
	case msg.Success:
		return r.SetStatus(sessionID, msg.Success.String())
	case msg.Warning:
		r.SetStatus(sessionID, msg.Warning.String())
		return r.SetWarning(sessionID, feedback.Content.(string))
	case msg.Cancelled:
		return r.SetStatus(sessionID, msg.Cancelled.String())
	case msg.Progress:
		switch feedback.Sender {
		case msg.InputConnector:
			return r.SetInputProgress(sessionID, feedback.Progress())
		case msg.OutputConnector:
			return r.SetOutputProgress(sessionID, feedback.Progress())
		}
	}
	return nil
}
func (r *Redis) Publish(sessionID string, feedback *msg.Feedback) error {
	f, err := feedback.Marshal()
	if err != nil {
		return err
	}
	return r.conn.Publish(getPubsubKey(sessionID), f).Err()
}
func (r *Redis) Subscribe(sessionID string) (*redis.PubSub, <-chan *msg.Feedback, error) {
	key := getPubsubKey(sessionID)
	pubsub := r.conn.Subscribe(key)
	ch := make(chan *msg.Feedback)
	go func() {
		for {
			m, err := pubsub.ReceiveMessage()
			if err != nil {
				close(ch)
				return
			}
			if m == nil || m.Payload == "" || m.Payload == "null" {
				ch <- nil
				continue
			}
			feedback := &msg.Feedback{}
			if err := feedback.Unmarshal([]byte(m.Payload)); err != nil {
				close(ch)
				return
			}
			ch <- feedback
		}
	}()
	return pubsub, ch, nil
}

func init() {
	backends.Add("redis", func(config map[string]interface{}) (msg.FeedbackBackend, error) {
		return New(config)
	})
}
