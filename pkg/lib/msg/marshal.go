package msg

import "encoding/json"

func (f *Feedback) Marshal() ([]byte, error) {
	return json.Marshal(f)
}

func (f *Feedback) Unmarshal(data []byte) error {
	return json.Unmarshal(data, f)
}
