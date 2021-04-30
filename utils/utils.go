package utils

import (
	"encoding/json"

	"github.com/imdario/mergo"
	"github.com/mitchellh/mapstructure"
)

type Field struct {
	Name      string
	Type      string
	TypeName  string
	Precision int64
	Scale     int64
	Nullable  bool
	Length    int64
}

func ParseOptions(result interface{}, options map[string]interface{}) {
	cfg := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   result,
		TagName:  "yaml",
	}
	decoder, _ := mapstructure.NewDecoder(cfg)
	decoder.Decode(options)
}

func Merge(dst, src interface{}, opts ...func(*mergo.Config)) error {
	return mergo.Merge(dst, src, opts...)
}

func MergeOptions(base, extend map[string]interface{}) (map[string]interface{}, error) {
	target := make(map[string]interface{})
	for key, value := range base {
		target[key] = value
	}
	err := Merge(&target, extend)
	if err != nil {
		return nil, err
	}
	return target, nil
}

func JSONLoadString(s string) (map[string]interface{}, error) {
	var m map[string]interface{} = nil
	if s != "" {
		m = make(map[string]interface{})
		err := json.Unmarshal([]byte(s), &m)
		if err != nil {
			return nil, err
		}
	}
	return m, nil
}
