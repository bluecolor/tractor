package tasks

import (
	"context"
	"encoding/json"

	"github.com/bluecolor/tractor/pkg/lib/runner"
	"github.com/hibiken/asynq"
)

func HandleExtractionTask(ctx context.Context, t *asynq.Task) error {
	p := ExtractionPayload{}
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}
	r, err := runner.New(ctx, p.SourceConnection, p.TargetConnection)
	if err != nil {
		return err
	}
	return r.Run(p.Params)
}
