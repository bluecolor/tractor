package tasks

import (
	"context"
	"encoding/json"

	"github.com/bluecolor/tractor/pkg/lib/msg"
	"github.com/bluecolor/tractor/pkg/lib/runner"
	"github.com/bluecolor/tractor/pkg/lib/types"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

func HandleExtractionTask(ctx context.Context, t *asynq.Task) error {
	s := types.Session{}
	if err := json.Unmarshal(t.Payload(), &s); err != nil {
		return err
	}
	log.Info().Msgf("Running extraction %s", s.Extraction.Name)

	backends := ctx.Value("backends").([]msg.FeedBackend)
	options := []runner.Option{}
	for _, backend := range backends {
		options = append(options, runner.WithFeedbackBackendOption(backend))
	}
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
