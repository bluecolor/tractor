package bridge

import (
	"github.com/bluecolor/tractor/pkg/lib/params"
	"github.com/bluecolor/tractor/pkg/models"
)

type Session struct {
	model models.Session
}

func NewSession(model models.Session) (s *Session) {
	return &Session{model: model}
}

func (s *Session) Session() (params.SessionParams, error) {
	e := NewExtraction(s.model.Extraction)
	session, err := e.Session()
	if err != nil {
		return nil, err
	}
	return session.WithSessionID(s.model.ID), nil
}
