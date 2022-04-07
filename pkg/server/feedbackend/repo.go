package feedbackend

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bluecolor/tractor/pkg/lib/msg"
	"github.com/bluecolor/tractor/pkg/models"
)

func (h *Handler) UpdateRepo(feed *msg.Feed) error {
	if !feed.IsSessionStatus() {
		return nil
	}
	status := strings.ToLower(feed.Type.String())
	sid, err := strconv.Atoi(feed.SessionID)
	if err != nil {
		return err
	}
	session := models.Session{}
	if err := h.repo.First(&session, sid).Error; err != nil {
		return err
	}
	if session.Status != status {
		session.Status = status
		h.updateSessionWithCache(&session)
		session.Log = fmt.Sprintf("%v", feed.Content)
		if err := h.repo.Save(&session).Error; err != nil {
			return err
		}
	}
	return nil
}

func (h *Handler) updateSessionWithCache(session *models.Session) error {
	return nil
}
