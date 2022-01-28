package providers

import (
	"github.com/bluecolor/tractor/pkg/lib/meta"
)

type Provider interface {
	FindDatasets(pattern string) ([]meta.Dataset, error)
}
