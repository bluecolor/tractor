package bridge

import (
	"encoding/json"

	types "github.com/bluecolor/tractor/pkg/lib/types"

	"gorm.io/datatypes"
)

func GetConfig(c datatypes.JSON) (types.Config, error) {
	if c == nil {
		return types.Config{}, nil
	}
	var config types.Config
	if err := json.Unmarshal(c, &config); err != nil {
		return nil, err
	}
	return config, nil
}
