package msg

import "encoding/json"

func (f *Feed) Marshal() ([]byte, error) {
	return json.Marshal(f)
}

func (f *Feed) Unmarshal(data []byte) error {
	return json.Unmarshal(data, f)
}
