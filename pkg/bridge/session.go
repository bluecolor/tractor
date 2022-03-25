package bridge

import (
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
