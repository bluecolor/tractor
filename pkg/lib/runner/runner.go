package runner

import "github.com/bluecolor/tractor/pkg/lib/meta"

type Runner struct {
	Params meta.ExtParams
}

func New(p meta.ExtParams) *Runner {
	return &Runner{Params: p}
}
