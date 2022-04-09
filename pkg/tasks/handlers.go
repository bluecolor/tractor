package tasks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bluecolor/tractor/pkg/lib/runner"
	"github.com/bluecolor/tractor/pkg/lib/types"
	"github.com/bluecolor/tractor/pkg/tasks/feedproc"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

func HandleExtractionTask(ctx context.Context, t *asynq.Task) error {
	s := types.Session{}
	if err := json.Unmarshal(t.Payload(), &s); err != nil {
		return err
	}
	log.Debug().Msgf("Running extraction %s", s.Extraction.Name)

	if ctx.Value("feed.processor") == nil {
		return fmt.Errorf("feed.processor is not set")
	}
	options := []runner.Option{
		runner.WithFeedClientOption(ctx.Value("feed.processor").(*feedproc.FeedProcessor)),
	}
	options = append(options)
	r, err := runner.New(ctx, s, options...)
	if err != nil {
		log.Error().Err(err).Msg("failed to create runner")
		return err
	}
	if err := r.Run(); err != nil {
		log.Error().Err(err).Msgf("failed to run extraction %s", s.Extraction.Name)
		return err
	}
	return nil
}
