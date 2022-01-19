package utils

import (
	"encoding/json"
	"strconv"

	"github.com/bluecolor/tractor/lib/config"
	"github.com/bluecolor/tractor/lib/feed"
	"github.com/imdario/mergo"
	"github.com/mitchellh/mapstructure"
)

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

func ParseOptions(result interface{}, options map[string]interface{}) {
	cfg := &mapstructure.DecoderConfig{
		Metadata: nil,
		Result:   result,
		TagName:  "yaml",
	}
	decoder, _ := mapstructure.NewDecoder(cfg)
	decoder.Decode(options)
}

func MapRecord(catalog *config.Catalog, record *feed.Record) (err error) {
	if catalog == nil || catalog.Fields == nil {
		return
	}
	fieldMap := catalog.GetFieldMap()
	for k, v := range *record {
		if _, ok := fieldMap[k]; !ok {
			continue
		}
		if fieldMap[k].Type == "integer" {
			switch val := v.(type) {
			case int, int8, int16, int32, int64:
				(*record)[k] = val
			case string:
				(*record)[k], err = strconv.ParseInt(val, 10, 64)
				if err != nil {
					return err
				}
			default:
			}
		}
	}
	return
}
