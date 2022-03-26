package dummy

import (
	"reflect"
	"time"

	"github.com/bluecolor/tractor/pkg/lib/msg"
	"github.com/brianvoe/gofakeit/v6"
)

type bar struct {
	Name   string
	Number int
	Float  float32
}

type fakerecord struct {
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

func (c *DummyConnector) Generate() <-chan interface{} {
	channel := make(chan interface{}, c.config.FakeRecordCount)
	go func() {
		data := msg.Data{}
		for i := 0; i < c.config.FakeRecordCount; i++ {
			fake := fakerecord{}
			gofakeit.Struct(&fake)
			record := msg.Record{}
			v := reflect.ValueOf(fake)
			for i := 0; i < v.NumField(); i++ {
				record = append(record, v.Field(i).Interface())
			}
			data = append(data, record)
			channel <- data
			data = msg.Data{}
			time.Sleep(time.Millisecond * time.Duration(c.config.FakeRecordInterval))
		}
		close(channel)
	}()
	return channel
}
