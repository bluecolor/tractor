package pubsub

import (
	"strings"

	"github.com/bluecolor/tractor/pkg/lib/msg"
	"github.com/go-redis/redis/v7"
)

type PubSub struct {
	client *redis.Client
}

func New(addr string) (*PubSub, error) {
	return &PubSub{
		client: redis.NewClient(&redis.Options{
			Addr: addr,
		}),
	}, nil
}
func (ps *PubSub) Close() error {
	return ps.client.Close()
}
func (ps *PubSub) TrySetRunning(sessionID string) error {
	session, err := ps.GetSession(sessionID)
	if err != nil {
		return err
	}
	_, ok := session[StatusKey]
	if !ok {
		ps.SetStatus(sessionID, "running")
	}
	return nil
}
func (ps *PubSub) GetSession(sessionID string) (map[string]string, error) {
	key := getSessionKey(sessionID)
	return ps.client.HGetAll(key).Result()
}
func (ps *PubSub) SetStatus(sessionID string, status string) error {
	key := getSessionKey(sessionID)
	return ps.client.HSet(key, StatusKey, strings.ToLower(status)).Err()
}
func (ps *PubSub) SetInputProgress(sessionID string, progress int) error {
	key := getSessionKey(sessionID)
	return ps.client.HSet(key, InputProgressKey, progress).Err()
}
func (ps *PubSub) SetOutputProgress(sessionID string, progress int) error {
	key := getSessionKey(sessionID)
	return ps.client.HSet(key, OutputProgressKey, progress).Err()
}
func (ps *PubSub) SetError(sessionID string, err error) error {
	key := getSessionKey(sessionID)
	return ps.client.HSet(key, ErrorKey, err.Error()).Err()
}
func (ps *PubSub) SetWarning(sessionID string, warning string) error {
	key := getSessionKey(sessionID)
	return ps.client.HSet(key, WarningKey, warning).Err()
}
func (ps *PubSub) Process(sessionID string, feedback *msg.Feedback) error {
	if err := ps.TrySetRunning(sessionID); err != nil {
		return err
	}
	if err := ps.Publish(sessionID, feedback); err != nil {
		return err
	}
	switch feedback.Type {
	case msg.Error:
		ps.SetStatus(sessionID, msg.Error.String())
		ps.SetError(sessionID, feedback.Error())
	case msg.Success:
		return ps.SetStatus(sessionID, msg.Success.String())
	case msg.Warning:
		ps.SetStatus(sessionID, msg.Warning.String())
		return ps.SetWarning(sessionID, feedback.Content.(string))
	case msg.Cancelled:
		return ps.SetStatus(sessionID, msg.Cancelled.String())
	case msg.Progress:
		switch feedback.Sender {
		case msg.InputConnector:
			return ps.SetInputProgress(sessionID, feedback.Progress())
		case msg.OutputConnector:
			return ps.SetOutputProgress(sessionID, feedback.Progress())
		}
	}
	return nil
}
func (ps *PubSub) Publish(sessionID string, feedback *msg.Feedback) error {
	f, err := feedback.Marshal()
	if err != nil {
		return err
	}
	return ps.client.Publish(getPubsubKey(sessionID), f).Err()
}
func (ps *PubSub) Subscribe(sessionID string) (*redis.PubSub, <-chan *msg.Feedback, error) {
	key := getPubsubKey(sessionID)
	pubsub := ps.client.Subscribe(key)
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
