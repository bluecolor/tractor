package bridge

import (
	"github.com/bluecolor/tractor/pkg/lib/types"
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

func (s *Session) SessionParams() (types.SessionParams, error) {
	return s.ext.SessionParams()
}

func (s *Session) Connections() (input *types.Connection, output *types.Connection, err error) {
	return s.ext.Connections()
}
