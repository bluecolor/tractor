package tasks

import (
	"encoding/json"

	"github.com/bluecolor/tractor/pkg/bridge"
	"github.com/bluecolor/tractor/pkg/models"
	"github.com/hibiken/asynq"
)

const (
	TypeEmailSend     = "email:send"
	TypeExtractionRun = "extraction:run"
)

func NewExtractionRunTask(e models.Extraction) (*asynq.Task, error) {
	ext := bridge.NewExtraction(e)
	params, err := ext.ExtParams()
	if err != nil {
		return nil, err
	}
	inputc, outputc, err := ext.Connections()
	if err != nil {
		return nil, err
	}
	payload, err := json.Marshal(ExtractionPayload{
		SourceConnection: inputc,
		TargetConnection: outputc,
		Params:           params,
	})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeExtractionRun, payload), nil
}
