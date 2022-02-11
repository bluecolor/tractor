package bridge

import (
	"encoding/json"

	"github.com/bluecolor/tractor/pkg/lib/params"
	"gorm.io/datatypes"
)

func GetConfig(c datatypes.JSON) (params.Config, error) {
	if c == nil {
		return params.Config{}, nil
	}
	var config params.Config
	if err := json.Unmarshal(c, &config); err != nil {
		return nil, err
	}
	return config, nil
}
