package file

import (
	"encoding/json"
	"testing"

	"github.com/bluecolor/tractor/pkg/lib/connectors"
)

func TestNewFileConnector(t *testing.T) {

	config, _ := json.Marshal(
		FileConfig{
			ProviderType: "localfs",
			Provider:     map[string]interface{}{"path": "/tmp"},
		},
	)
	connector, err := NewFileConnector(connectors.ConnectorConfig(config))
	if err != nil {
		t.Error(err)
	}
	if connector == nil {
		t.Error("connector is nil")
	}
}
