package tasks

import (
	"context"
	"encoding/json"

	"github.com/bluecolor/tractor/pkg/lib/runner"
	"github.com/bluecolor/tractor/pkg/lib/types"
	"github.com/hibiken/asynq"
)

func HandleExtractionTask(ctx context.Context, t *asynq.Task) error {
	s := types.Session{}
	if err := json.Unmarshal(t.Payload(), &s); err != nil {
		return err
	}
	r, err := runner.New(ctx, s)
	if err != nil {
		return err
	}
	return r.Run()
}
