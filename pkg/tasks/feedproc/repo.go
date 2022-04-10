package feedproc

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bluecolor/tractor/pkg/lib/msg"
	"github.com/bluecolor/tractor/pkg/lib/types"
	"github.com/bluecolor/tractor/pkg/models"
	"github.com/rs/zerolog/log"
)

func (fp *FeedProcessor) UpdateRepo(feed *msg.Feed) error {
	if !feed.IsSessionStatus() {
		return nil
	}
	status := strings.ToLower(feed.Type.String())
	sid, err := strconv.Atoi(feed.SessionID)
	if err != nil {
		return err
	}
	session := models.Session{}
	if err := fp.repo.First(&session, sid).Error; err != nil {
		return err
	}
	if session.Status != status && status != types.StatusDone {
		session.Status = status
		fp.updateSessionWithCache(&session)
		session.Logs = fmt.Sprintf("%v", feed.Content)
		if status == types.StatusRunning {
			now := time.Now()
			session.StartedAt = &now
		} else if status == types.StatusSuccess || status == types.StatusError {
			now := time.Now()
			session.FinishedAt = &now
		}
		if err := fp.repo.Save(&session).Error; err != nil {
			return err
		}
	}
	return nil
}

func (fp *FeedProcessor) updateSessionWithCache(session *models.Session) error {
	sessionCache, err := fp.cache.HGetAll(ctx, getSessionKey(session.ID)).Result()
	if err != nil {
		log.Error().Err(err).Msg("failed to get session from cache")
		return err
	}
	if _, ok := sessionCache["input_progress"]; ok {
		progress, err := strconv.Atoi(sessionCache["input_progress"])
		if err != nil {
			return err
		}
		session.ReadCount = progress
	}
	if _, ok := sessionCache["output_progress"]; ok {
		progress, err := strconv.Atoi(sessionCache["output_progress"])
		if err != nil {
			return err
		}
		session.WriteCount = progress
	}
	return nil
}
