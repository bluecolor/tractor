package bridge

import (
	"testing"

	"github.com/bluecolor/tractor/pkg/lib/types"
	"github.com/bluecolor/tractor/pkg/models"
)

func TestNewSession(t *testing.T) {
	t.Parallel()
	tests := []struct {
		Model    models.Extraction
		Expected types.SessionParams
	}{
		{
			Model: models.Extraction{
				SourceDataset: &models.Dataset{Name: "source"},
				TargetDataset: &models.Dataset{Name: "target"},
				FieldMappings: []models.FieldMapping{
					{
						SourceField: &models.Field{Name: "source_field"},
						TargetField: &models.Field{Name: "target_field"},
					},
				},
			},
			Expected: types.SessionParams{
				types.InputDatasetKey: &types.Dataset{
					Name:   "source",
					Fields: []*types.Field{{Name: "source_field"}},
				},
				types.OutputDatasetKey: &types.Dataset{
					Name:   "target",
					Fields: []*types.Field{{Name: "target_field"}},
				},
				types.FieldMappingsKey: []types.FieldMapping{
					{
						SourceField: &types.Field{Name: "source_field"},
						TargetField: &types.Field{Name: "target_field"},
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			e := NewExtraction(&test.Model)
			p, err := e.SessionParams()
			if err != nil {
				t.Error(err)
			}
			if p.GetInputDataset().Name != test.Expected.GetInputDataset().Name {
				t.Fatalf("expected %s, got %s", test.Expected.GetInputDataset().Name, p.GetInputDataset().Name)
			}
			if len(p.GetInputDataset().Fields) != len(test.Expected.GetInputDataset().Fields) {
				t.Fatalf("expected %d, got %d", len(test.Expected.GetInputDataset().Fields), len(p.GetInputDataset().Fields))
			}
			if p.GetOutputDataset().Name != test.Expected.GetOutputDataset().Name {
				t.Fatalf("expected %s, got %s", test.Expected.GetOutputDataset().Name, p.GetOutputDataset().Name)
			}
			if len(p.GetOutputDataset().Fields) != len(test.Expected.GetOutputDataset().Fields) {
				t.Fatalf("expected %d, got %d", len(test.Expected.GetOutputDataset().Fields), len(p.GetOutputDataset().Fields))
			}
			if len(p.GetFMInputFields()) != len(test.Expected.GetFMInputFields()) {
				t.Fatalf("expected %d, got %d", len(test.Expected.GetFMInputFields()), len(p.GetFMInputFields()))
			}
			if len(p.GetFMOutputFields()) != len(test.Expected.GetFMOutputFields()) {
				t.Fatalf("expected %d, got %d", len(test.Expected.GetFMOutputFields()), len(p.GetFMOutputFields()))
			}
		})
	}
}
