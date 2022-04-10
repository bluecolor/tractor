package feedproc

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/bluecolor/tractor/pkg/lib/msg"
	"github.com/rs/zerolog/log"
)

var ctx = context.Background()

func getSessionKey[T string | int | uint](sessionID T) string {
	return fmt.Sprintf("tractor:session:%v", sessionID)
}
func getPubsubKey() string {
	return fmt.Sprintf("tractor:session:feeds")
}

func (fp *FeedProcessor) UpdateCache(feed *msg.Feed) error {
	session, err := fp.cache.HGetAll(ctx, getSessionKey(feed.SessionID)).Result()
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
	if err := fp.cache.HSet(ctx, getSessionKey(feed.SessionID), session).Err(); err != nil {
		log.Error().Err(err).Msg("failed to update cache")
		return err
	}
	return nil
}

func (fp *FeedProcessor) Publish(feed *msg.Feed) error {
	payload, err := feed.Marshal()
	if err != nil {
		return err
	}
	key := getPubsubKey()
	if err := fp.cache.Publish(ctx, key, payload).Err(); err != nil {
		log.Error().Err(err).Msg("failed to publish to " + key)
		return err
	}
	return nil
}
