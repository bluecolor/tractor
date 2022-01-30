package file

import (
	"testing"

	"github.com/bluecolor/tractor/pkg/lib/connectors"
)

func TestNewFileConnector(t *testing.T) {
	configs := []connectors.ConnectorConfig{
		{
			"storageType": "fs",
			"format":      "csv",
			"storageConfig": map[string]interface{}{
				"url": "fs:///tmp/",
			},
		},
	}
	for _, config := range configs {
		if _, err := New(config); err != nil {
			t.Error(err)
		}
	}
}
