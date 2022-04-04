package repository

import (
	"strconv"

	"github.com/bluecolor/tractor/pkg/lib/msg"
	"github.com/bluecolor/tractor/pkg/models"
	"github.com/bluecolor/tractor/pkg/repo"
)

type Repository struct {
	repo *repo.Repository
}

func New(r *repo.Repository) *Repository {
	return &Repository{
		repo: r,
	}
}

func (r *Repository) Process(sessionID string, feedback *msg.Feedback) error {
	if err := r.setStatus(sessionID, feedback); err != nil {
		return err
	}
	return nil
}

func (r *Repository) setStatus(sessionID string, feedback *msg.Feedback) error {
	var status string
	if feedback.IsSessionRunning() {
		status = models.SessionStatusRunning
	} else if feedback.IsSessionSuccess() {
		status = models.SessionStatusSuccess
	} else if feedback.IsSessionError() {
		status = models.SessionStatusError
	}
	if status == "" {
		return nil
	}
	sid, err := strconv.Atoi(sessionID)
	if err != nil {
		return err
	}
	session := models.Session{}
	if err := r.repo.First(&session, sid).Error; err != nil {
		return err
	}
	if session.Status != status {
		session.Status = status
		if err := r.repo.Save(&session).Error; err != nil {
			return err
		}
	}
	return nil
}
