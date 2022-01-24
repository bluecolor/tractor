package ext

type Field struct {
	Name   string
	Type   string
	Config map[string]string
}

type SourceField struct {
	Field
	Expression string
}
type TargetField struct {
	Field
}

type Dataset struct {
	Name   string
	Fields []Field
	Config map[string]string
}

type FieldMapping struct {
	SourceField SourceField
	TargetField TargetField
}

type Extraction struct {
	SourceDataset  Dataset
	TargetDataset  Dataset
	ExtractionMode string
}
