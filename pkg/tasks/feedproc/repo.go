package feedproc

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bluecolor/tractor/pkg/lib/msg"
	"github.com/bluecolor/tractor/pkg/models"
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
	if session.Status != status {
		session.Status = status
		fp.updateSessionWithCache(&session)
		session.Log = fmt.Sprintf("%v", feed.Content)
		if err := fp.repo.Save(&session).Error; err != nil {
			return err
		}
	}
	return nil
}

func (fp *FeedProcessor) updateSessionWithCache(session *models.Session) error {
	return nil
}
