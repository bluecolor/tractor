package tasks

import (
	"encoding/json"

	"github.com/bluecolor/tractor/pkg/bridge"
	"github.com/bluecolor/tractor/pkg/models"
	"github.com/hibiken/asynq"
)

const (
	TypeEmailSend  = "email:send"
	TypeSessionRun = "session:run"
)

func NewSessionRunTask(s *models.Session) (*asynq.Task, error) {
	session := bridge.NewSession(s)
	params, err := session.SessionParams()
	if err != nil {
		return nil, err
	}
	inputc, outputc, err := session.Connections()
	if err != nil {
		return nil, err
	}
	payload, err := json.Marshal(ExtractionPayload{
		SourceConnection: inputc,
		TargetConnection: outputc,
		Dataset:          data,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeSessionRun, payload), nil
}
