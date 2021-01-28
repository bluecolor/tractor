package api

import (
	"plugin"
	"reflect"
	"sync"
)

// PluginType ...
type PluginType int

// MetadataType ...
type MetadataType int

// NodeType ...
type NodeType int

// FeedType ...
type FeedType int

//Config either input our ooutput configuration given by the user
//in mappings.yml file
type Config map[interface{}]interface{}

const (
	// InputPlugin ...
	InputPlugin PluginType = iota
	// OutputPlugin ...
	OutputPlugin
)

const (
	// FieldsMetadata ...
	FieldsMetadata MetadataType = iota
)

const (
	// InputNode ...
	InputNode NodeType = iota
	// OutputNode ...
	OutputNode
	// Master ...
	Master
)

const (
	// ReadCountFeed ...
	ReadCountFeed FeedType = iota
	// WriteCountFeed ...
	WriteCountFeed
)

// TractorPlugin ...
type TractorPlugin struct {
	Plugin *plugin.Plugin
	Run    func(*sync.WaitGroup, []byte, *Wire) error
}

// DecimalSize ...
type DecimalSize struct {
	Precision int64
	Scale     int64
	Ok        bool
}

// Nullable ...
type Nullable struct {
	Nullable bool
	Ok       bool
}

// Length ...
type Length struct {
	Length int64
	Ok     bool
}

// Field ...
type Field struct {
	Name        string
	Type        reflect.Type
	DecimalSize DecimalSize
	Nullable    Nullable
	Length      Length
}

// Metadata ...
type Metadata struct {
	Type    MetadataType
	Content interface{}
}

// Data ...
type Data struct {
	Content [][]interface{}
}

// Feed ...
type Feed struct {
	Type    FeedType
	Sender  NodeType
	Content interface{}
}

// NewReadCountFeed ...
func NewReadCountFeed(count int) *Feed {
	return &Feed{
		Type:    ReadCountFeed,
		Sender:  InputNode,
		Content: count,
	}
}

// NewWriteCountFeed ...
func NewWriteCountFeed(count int) *Feed {
	return &Feed{
		Type:    WriteCountFeed,
		Sender:  OutputNode,
		Content: count,
	}
}

// Wire ...
type Wire struct {
	Feed     chan *Feed
	Metadata chan *Metadata
	Data     chan *Data
}

// CloseFeedChannel ...
func (w *Wire) CloseFeedChannel() {
	close(w.Feed)
}

// CloseDataChannel ...
func (w *Wire) CloseDataChannel() {
	close(w.Data)
}

// CloseMetadataChannel ...
func (w *Wire) CloseMetadataChannel() {
	close(w.Metadata)
}
