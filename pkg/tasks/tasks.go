package tasks

import (
	"encoding/json"

	"github.com/bluecolor/tractor/pkg/lib/types"
	"github.com/hibiken/asynq"
)

const (
	TypeEmailSend  = "email:send"
	TypeSessionRun = "session:run"
)

func NewSessionRunTask(s *types.Session) (*asynq.Task, error) {
	payload, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeSessionRun, payload), nil
}
