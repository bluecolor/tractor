package test

import (
	"encoding/json"
	"time"

	"github.com/bluecolor/tractor/pkg/lib/msg"
	"github.com/bluecolor/tractor/pkg/lib/types"
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

func GetSession() types.SessionParams {
	inChannel := make(chan interface{}, 1000)
	outChannel := make(chan interface{}, 1000)
	inputDataset := types.Dataset{
		Name: "test_input",
		Fields: []*types.Field{
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
	outputDataset := types.Dataset{
		Name: "test_output",
		Fields: []*types.Field{
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
	fm := []types.FieldMapping{
		{
			SourceField: &types.Field{Name: "name"},
			TargetField: &types.Field{Name: "name"},
		},
		{
			SourceField: &types.Field{Name: "randstr"},
			TargetField: &types.Field{Name: "randstr"},
		},
		{
			SourceField: &types.Field{Name: "number"},
			TargetField: &types.Field{Name: "number"},
		},
		{
			SourceField: &types.Field{Name: "regex"},
			TargetField: &types.Field{Name: "regex"},
		},
		{
			SourceField: &types.Field{Name: "map"},
			TargetField: &types.Field{Name: "map"},
		},
		{
			SourceField: &types.Field{Name: "array"},
			TargetField: &types.Field{Name: "array"},
		},
		{
			SourceField: &types.Field{Name: "bar"},
			TargetField: &types.Field{Name: "bar"},
		},
		{
			SourceField: &types.Field{Name: "created"},
			TargetField: &types.Field{Name: "created"},
		},
		{
			SourceField: &types.Field{Name: "created_format"},
			TargetField: &types.Field{Name: "created_format"},
		},
	}
	return types.SessionParams{}.
		WithInputDataset(&inputDataset).
		WithFieldMappings(fm).
		WithOutputDataset(&outputDataset)
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

func getOneRecord() (record msg.Record, err error) {
	fake := testrecord{}
	gofakeit.Struct(&fake)

	r, err := json.Marshal(fake)
	if err != nil {
		return nil, err
	}
	record = msg.Record{}
	if err = json.Unmarshal(r, &record); err != nil {
		return nil, err
	}
	return record, nil
}

func GenerateTestDataWithDuration(rc int, ch chan<- interface{}, dur time.Duration) (err error) {
	for i := 0; i < rc; i++ {
		time.Sleep(dur / time.Duration(rc))
		record, err := getOneRecord()
		if err != nil {
			return err
		}
		ch <- record
	}
	return
}
