package tasks

import "github.com/bluecolor/tractor/pkg/lib/params"

type ExtractionPayload struct {
	SourceConnection *params.Connection `json:"source_connection"`
	TargetConnection *params.Connection `json:"target_connection"`
	Params           params.ExtParams   `json:"params"`
}
