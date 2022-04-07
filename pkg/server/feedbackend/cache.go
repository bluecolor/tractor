package feedbackend

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bluecolor/tractor/pkg/lib/msg"
)

func (h *Handler) UpdateCache(feed *msg.Feed) error {
	session, err := h.cache.HGetAll(getSessionKey(feed.SessionID)).Result()
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
	return h.cache.HMSet(getSessionKey(feed.SessionID), session).Err()
}

func (h *Handler) Publish(feed *msg.Feed) error {
	payload, err := feed.Marshal()
	if err != nil {
		return err
	}
	return h.cache.Publish(getPubsubKey(), payload).Err()
}

// func (f *FeedBackend) Subscribe() (*redis.PubSub, <-chan *msg.Feed, error) {
// 	pubsub := f.cache.Subscribe(getPubsubKey())
// 	ch := make(chan *msg.Feed)
// 	go func() {
// 		for {
// 			m, err := pubsub.ReceiveMessage()
// 			if err != nil {
// 				close(ch)
// 				return
// 			}
// 			if m == nil || m.Payload == "" || m.Payload == "null" {
// 				ch <- nil
// 				continue
// 			}
// 			feed := &msg.Feed{}
// 			if err := feed.Unmarshal([]byte(m.Payload)); err != nil {
// 				close(ch)
// 				return
// 			}
// 			ch <- feed
// 		}
// 	}()
// 	return pubsub, ch, nil
// }
