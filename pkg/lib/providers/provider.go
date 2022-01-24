package providers

import "github.com/bluecolor/tractor/pkg/lib/cat/meta"

type Provider interface {
	FindDatasets(pattern string) ([]meta.Dataset, error)
}
