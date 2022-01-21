package providers

import "github.com/bluecolor/tractor/pkg/models"

type Provider interface {
	FetchDatasetsWithPattern(pattern string) ([]models.Dataset, error)
}
