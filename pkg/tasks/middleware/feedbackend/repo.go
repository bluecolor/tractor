package feedbackend

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bluecolor/tractor/pkg/lib/msg"
	"github.com/bluecolor/tractor/pkg/models"
)

func (f *FeedBackend) UpdateRepo(sessionID string, feed *msg.Feed) error {
	if !feed.IsSessionStatus() {
		return nil
	}
	status := strings.ToLower(feed.Type.String())
	sid, err := strconv.Atoi(sessionID)
	if err != nil {
		return err
	}
	session := models.Session{}
	if err := f.repo.First(&session, sid).Error; err != nil {
		return err
	}
	if session.Status != status {
		session.Status = status
		f.updateSessionWithCache(&session)
		session.Log = fmt.Sprintf("%v", feed.Content)
		if err := f.repo.Save(&session).Error; err != nil {
			return err
		}
	}
	return nil
}

func (f *FeedBackend) updateSessionWithCache(session *models.Session) error {
	return nil
}
