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

func (s *Session) SessionParams() (params.SessionParams, error) {
	return s.ext.SessionParams()
}

func (s *Session) Connections() (input *params.Connection, output *params.Connection, err error) {
	return s.ext.Connections()
}
