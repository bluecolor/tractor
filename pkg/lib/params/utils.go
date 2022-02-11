package params

import (
	"github.com/bluecolor/tractor/pkg/lib/msg"
	"github.com/rs/zerolog/log"
)

func ToOutputData(input []msg.Record, p ExtParams) ([]msg.Record, error) {
	dataset := p.GetOutputDataset()
	output := make([]msg.Record, len(input))
	for i, r := range input {
		record := make(msg.Record, len(dataset.Fields))
		for _, f := range dataset.Fields {
			sourceFieldName := p.GetSourceFieldNameByTargetFieldName(f.Name)
			v, ok := r[sourceFieldName]
			record[f.Name] = v
			if !ok {
				log.Debug().Msgf("source field matching to target %s not found in record %d", f.Name, i)
			}
		}
		output[i] = record
	}
	return output, nil
}
