package bridge

import (
	"github.com/bluecolor/tractor/pkg/lib/params"
	"github.com/bluecolor/tractor/pkg/models"
)

type Session struct {
	model *models.Session
	ext   *Extraction
}

func NewSession(model *models.Session) (s *Session) {
	return &Session{
		model: model,
		ext:   NewExtraction(s.model.Extraction),
	}
}

func (s *Session) SessionParamsWithID() (params.SessionParams, error) {
	session, err := s.ext.SessionParams()
	if err != nil {
		return nil, err
	}
	return session.WithSessionID(s.model.ID), nil
}

func (s *Session) Connections() (input *params.Connection, output *params.Connection, err error) {
	return s.ext.Connections()
}
