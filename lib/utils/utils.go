package utils

import "encoding/json"

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
