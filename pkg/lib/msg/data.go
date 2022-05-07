package msg

type (
	Record []interface{}
	Data   []Record
)

func NewData(data interface{}, args ...interface{}) Data {
	var content []Record
	switch val := data.(type) {
	case []Record:
		content = val
	case Record:
		content = []Record{val}
	case Data:
		content = val
	default:
		return nil
	}
	return content
}

func (d Data) Count() int {
	return len(d)
}
