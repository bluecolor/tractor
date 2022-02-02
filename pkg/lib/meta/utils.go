package meta

import (
	"github.com/bluecolor/tractor/pkg/lib/feeds"
	"github.com/rs/zerolog/log"
)

func ToOutputData(inputData feeds.Data, p ExtParams) (feeds.Data, error) {
	dataset := p.GetOutputDataset()
	outputData := make(feeds.Data, len(inputData))
	for i, r := range inputData {
		record := make(feeds.Record, len(dataset.Fields))
		for _, f := range dataset.Fields {
			sourceFieldName := p.GetSourceFieldNameByTargetFieldName(f.Name)
			v, ok := r[sourceFieldName]
			record[f.Name] = v
			if !ok {
				log.Debug().Msgf("source field matching to target %s not found in record %d", f.Name, i)
			}
		}
		outputData[i] = record
	}
	return outputData, nil
}
