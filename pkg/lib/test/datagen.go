package test

import (
	"encoding/json"
	"time"

	"github.com/bluecolor/tractor/pkg/lib/meta"
	"github.com/bluecolor/tractor/pkg/lib/msg"
	"github.com/brianvoe/gofakeit/v6"
)

type bar struct {
	Name   string
	Number int
	Float  float32
}
type testrecord struct {
	Str           string         `json:"str"`
	Int           int            `json:"int"`
	Pointer       *int           `json:"pointer"`
	Name          string         `fake:"{firstname}" json:"name"`
	Sentence      string         `fake:"{sentence:3}" json:"sentence"`
	RandStr       string         `fake:"{randomstring:[hello,world]}" json:"randstr"`
	Number        string         `fake:"{number:1,10}" json:"number"`
	Regex         string         `fake:"{regex:[abcdef]{5}}" json:"regex"`
	Map           map[string]int `fakesize:"2" json:"map"`
	Array         []string       `fakesize:"2" json:"array"`
	Bar           bar            `json:"bar"`
	Skip          *string        `fake:"skip" json:"-"`
	Created       time.Time      `json:"created"`
	CreatedFormat time.Time      `fake:"{year}-{month}-{day}" format:"2006-01-02" json:"created_format"`
}

func GetExtParams() meta.ExtParams {
	inChannel := make(chan interface{}, 1000)
	outChannel := make(chan interface{}, 1000)
	inputDataset := meta.Dataset{
		Name: "test_input",
		Fields: []meta.Field{
			{Name: "str", Type: "string"},
			{Name: "int", Type: "int"},
			{Name: "pointer", Type: "int"},
			{Name: "name", Type: "string"},
			{Name: "sentence", Type: "string"},
			{Name: "randstr", Type: "string"},
			{Name: "number", Type: "string"},
			{Name: "regex", Type: "string"},
			{Name: "map", Type: "map"},
			{Name: "array", Type: "array"},
			{Name: "bar", Type: "struct"},
			{Name: "created", Type: "time"},
			{Name: "created_format", Type: "time"},
		},
		Config: map[string]interface{}{
			"channel": inChannel,
		},
	}
	outputDataset := meta.Dataset{
		Name: "test_output",
		Fields: []meta.Field{
			{Name: "str", Type: "string"},
			{Name: "int", Type: "int"},
			{Name: "pointer", Type: "int"},
			{Name: "name", Type: "string"},
			{Name: "sentence", Type: "string"},
			{Name: "randstr", Type: "string"},
			{Name: "number", Type: "string"},
			{Name: "regex", Type: "string"},
			{Name: "map", Type: "map"},
			{Name: "array", Type: "array"},
			{Name: "bar", Type: "struct"},
			{Name: "created", Type: "time"},
			{Name: "created_format", Type: "time"},
		},
		Config: map[string]interface{}{
			"channel": outChannel,
		},
	}
	fm := []meta.FieldMapping{
		{
			SourceField: meta.Field{Name: "name"},
			TargetField: meta.Field{Name: "name"},
		},
		{
			SourceField: meta.Field{Name: "randstr"},
			TargetField: meta.Field{Name: "randstr"},
		},
		{
			SourceField: meta.Field{Name: "number"},
			TargetField: meta.Field{Name: "number"},
		},
		{
			SourceField: meta.Field{Name: "regex"},
			TargetField: meta.Field{Name: "regex"},
		},
		{
			SourceField: meta.Field{Name: "map"},
			TargetField: meta.Field{Name: "map"},
		},
		{
			SourceField: meta.Field{Name: "array"},
			TargetField: meta.Field{Name: "array"},
		},
		{
			SourceField: meta.Field{Name: "bar"},
			TargetField: meta.Field{Name: "bar"},
		},
		{
			SourceField: meta.Field{Name: "created"},
			TargetField: meta.Field{Name: "created"},
		},
		{
			SourceField: meta.Field{Name: "created_format"},
			TargetField: meta.Field{Name: "created_format"},
		},
	}
	return meta.ExtParams{}.
		WithInputDataset(inputDataset).
		WithFieldMappings(fm).
		WithOutputDataset(outputDataset)
}
func GenerateTestData(recordCount int, ch chan<- interface{}) (err error) {
	data := []msg.Record{}
	for i := 0; i < recordCount; i++ {
		fake := testrecord{}
		gofakeit.Struct(&fake)
		r, err := json.Marshal(fake)
		if err != nil {
			return err
		}
		record := msg.Record{}
		if err = json.Unmarshal(r, &record); err != nil {
			return err
		}
		data = append(data, record)
		ch <- data
		data = []msg.Record{}
	}
	return
}
