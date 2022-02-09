package meta

type ExtractionMode int

const (
	ExtractionModeCreate ExtractionMode = iota
	ExtractionModeInsert
	ExtractionModeUpdate
	ExtractionModeDelete
)

func (m ExtractionMode) String() string {
	switch m {
	case ExtractionModeCreate:
		return "create"
	case ExtractionModeInsert:
		return "insert"
	case ExtractionModeUpdate:
		return "update"
	case ExtractionModeDelete:
		return "delete"
	}
	return "create"
}
func ExtractionModeFromString(mode string) ExtractionMode {
	var m ExtractionMode
	return m.FromString(mode)
}
func (m ExtractionMode) FromString(mode string) ExtractionMode {
	switch mode {
	case "create":
		return ExtractionModeCreate
	case "insert":
		return ExtractionModeInsert
	case "update":
		return ExtractionModeUpdate
	case "delete":
		return ExtractionModeDelete
	}
	return ExtractionModeCreate
}

type Config map[string]interface{}

func (c Config) GetString(key string, def string) string {
	if v, ok := c[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return def
}
func (c Config) GetInt(key string, def int) int {
	if v, ok := c[key]; ok {
		if i, ok := v.(int); ok {
			return i
		}
	}
	return def
}
func (c Config) GetStringArray(key string, def []string) []string {
	if v, ok := c[key]; ok {
		if s, ok := v.([]string); ok {
			return s
		}
	}
	return def
}
func (c Config) GetBool(key string, def bool) bool {
	if v, ok := c[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return def
}
func (c Config) GetChannel(key string) chan interface{} {
	if v, ok := c[key]; ok {
		if ch, ok := v.(chan interface{}); ok {
			return ch
		}
	}
	return nil
}

type Dataset struct {
	Name   string  `json:"name"`
	Config Config  `json:"config"`
	Fields []Field `json:"fields"`
}
type Field struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Expression string `json:"expression"`
	Config     Config `json:"config"`
}

func (f *Field) GetExpressionOrName() string {
	if f.Expression != "" {
		return f.Expression
	}
	return f.Name
}

type FieldMapping struct {
	SourceField Field  `json:"sourceField"`
	TargetField Field  `json:"targetField"`
	Config      Config `json:"config"`
}

type Connection struct {
	Name           string `json:"name"`
	ConnectionType string `json:"connectionType"`
	Config         Config `json:"config"`
	AsSource       bool   `json:"asSource"`
	AsTarget       bool   `json:"asTarget"`
}
