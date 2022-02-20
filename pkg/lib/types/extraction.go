package types

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
