package feedbackend

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bluecolor/tractor/pkg/lib/msg"
	"github.com/go-redis/redis/v7"
)

func (f *FeedBackend) UpdateCache(sessionID string, feed *msg.Feed) error {
	session, err := f.GetSessionCache(sessionID)
	if err != nil {
		return err
	}
	if feed.IsSessionStatus() {
		session["status"] = strings.ToLower(feed.Type.String())
		session["log"] = fmt.Sprintf("%v", feed.Content)
	}
	if feed.IsProgress() {
		var inputProgress, outputProgress int = 0, 0
		if _, ok := session["input_progress"]; ok {
			progress, err := strconv.Atoi(session["input_progress"])
			if err != nil {
				return err
			}
			inputProgress += progress
		}
		if _, ok := session["output_progress"]; ok {
			progress, err := strconv.Atoi(session["output_progress"])
			if err != nil {
				return err
			}
			outputProgress += progress
		}
		inputProgress += feed.InputProgress()
		outputProgress += feed.OutputProgress()
		session["input_progress"] = strconv.Itoa(inputProgress)
		session["output_progress"] = strconv.Itoa(outputProgress)
	}
	return f.cache.HMSet(getSessionKey(sessionID), session).Err()
}

func (f *FeedBackend) GetSessionCache(sessionID string) (map[string]string, error) {
	key := getSessionKey(sessionID)
	return f.cache.HGetAll(key).Result()
}

func (f *FeedBackend) Publish(sessionID string, feed *msg.Feed) error {
	payload, err := feed.Marshal()
	if err != nil {
		return err
	}
	return f.cache.Publish(getPubsubKey(sessionID), payload).Err()
}

func (f *FeedBackend) Subscribe(sessionID string) (*redis.PubSub, <-chan *msg.Feed, error) {
	pubsub := f.cache.Subscribe(getPubsubKey(sessionID))
	ch := make(chan *msg.Feed)
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
			feedback := &msg.Feed{}
			if err := feedback.Unmarshal([]byte(m.Payload)); err != nil {
				close(ch)
				return
			}
			ch <- feedback
		}
	}()
	return pubsub, ch, nil
}
