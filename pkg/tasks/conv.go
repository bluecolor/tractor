package tasks

import (
	"encoding/json"
	"fmt"

	"github.com/bluecolor/tractor/pkg/lib/types"
	"github.com/bluecolor/tractor/pkg/models"
	"gorm.io/datatypes"
)

func getConfig(c datatypes.JSON) (types.Config, error) {
	if c == nil {
		return types.Config{}, nil
	}
	var config types.Config
	if err := json.Unmarshal(c, &config); err != nil {
		return nil, err
	}
	return config, nil
}

func getFields(fields []*models.Field) (output []*types.Field, err error) {
	output = make([]*types.Field, len(fields))
	for i, f := range fields {
		config, err := getConfig(f.Config)
		if err != nil {
			return nil, err
		}
		output[i] = &types.Field{
			Name:   f.Name,
			Type:   types.FieldTypeFromString(f.Type),
			Config: config,
		}
	}
	return output, nil
}
func getConnection(connection *models.Connection) (*types.Connection, error) {
	config, err := getConfig(connection.Config)
	if err != nil {
		return nil, err
	}
	return &types.Connection{
		Name:           connection.Name,
		Config:         config,
		ConnectionType: connection.ConnectionType.Code,
	}, nil
}
func getDataset(dataset *models.Dataset) (*types.Dataset, error) {
	config, err := getConfig(dataset.Config)
	if err != nil {
		return nil, err
	}
	fields, err := getFields(dataset.Fields)
	if err != nil {
		return nil, err
	}
	connection, err := getConnection(dataset.Connection)
	return &types.Dataset{
		Name:       dataset.Name,
		Fields:     fields,
		Config:     config,
		Connection: connection,
	}, nil
}
func getExtraction(extraction *models.Extraction) (*types.Extraction, error) {
	config, err := getConfig(extraction.Config)
	if err != nil {
		return nil, err
	}
	if extraction.SourceDataset == nil {
		return nil, fmt.Errorf("extraction.SourceDataset is nil")
	}
	sourceDataset, err := getDataset(extraction.SourceDataset)
	if err != nil {
		return nil, err
	}
	if extraction.TargetDataset == nil {
		return nil, fmt.Errorf("extraction.TargetDataset is nil")
	}
	targetDataset, err := getDataset(extraction.TargetDataset)
	if err != nil {
		return nil, err
	}
	return &types.Extraction{
		SourceDataset: sourceDataset,
		TargetDataset: targetDataset,
		Config:        config,
	}, nil
}
func GetSession(s *models.Session) (*types.Session, error) {
	config, err := getConfig(s.Config)
	if err != nil {
		return nil, err
	}
	extraction, err := getExtraction(s.Extraction)
	if err != nil {
		return nil, err
	}
	return &types.Session{
		ID:         fmt.Sprintf("%d", s.ID),
		Config:     config,
		Extraction: extraction,
	}, nil
}
