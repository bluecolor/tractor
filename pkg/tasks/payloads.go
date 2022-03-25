package tasks

import "github.com/bluecolor/tractor/pkg/lib/types"

type ExtractionPayload struct {
	SourceConnection *types.Connection `json:"source_connection"`
	TargetConnection *types.Connection `json:"target_connection"`
	Dataset          types.Dataset     `json:"dataset"`
}
