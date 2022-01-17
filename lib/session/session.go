package session

import (
	"fmt"
	"sync"
	"time"

	"github.com/bluecolor/tractor/lib/config"
	"github.com/bluecolor/tractor/lib/feed"
	"github.com/bluecolor/tractor/lib/plugins"
	"github.com/bluecolor/tractor/lib/plugins/inputs"
	"github.com/bluecolor/tractor/lib/plugins/outputs"
	"github.com/bluecolor/tractor/lib/utils"
	"github.com/bluecolor/tractor/lib/wire"
)

type Session struct {
	StartTime    time.Time
	Mapping      *config.Mapping
	Progress     *feed.ProgressFeed
	InputPlugin  inputs.InputPlugin
	OutputPlugin outputs.OutputPlugin
	Wire         wire.Wire
}

func (s *Session) Start() (*sync.WaitGroup, error) {
	var wg sync.WaitGroup
	go func(wg *sync.WaitGroup) {
		s.InputPlugin.Read(s.Wire)
		wg.Done()
		s.Wire.CloseData()
	}(&wg)
	wg.Add(1)

	go func(wg *sync.WaitGroup) {
		s.OutputPlugin.Write(s.Wire)
		wg.Done()
		s.Wire.CloseFeed()
	}(&wg)
	wg.Add(1)

	return &wg, nil
}
func (s *Session) CreateInputPlugin(params string) (err error) {
	if creator, ok := inputs.Inputs[s.Mapping.Input.Plugin]; ok {
		var paramsMap map[string]interface{} = nil
		if params != "" {
			paramsMap, err = utils.JSONLoadString(params)
			if err != nil {
				return err
			}
		}
		s.InputPlugin, err = creator(s.Mapping.Input.Config, s.Mapping.Input.Catalog, paramsMap)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("plugin %s not found", s.Mapping.Input.Plugin)
	}
	return nil
}
func (s *Session) CreateOutputPlugin(params string) (err error) {
	if creator, ok := outputs.Outputs[s.Mapping.Output.Plugin]; ok {
		var paramsMap map[string]interface{} = nil
		if params != "" {
			paramsMap, err = utils.JSONLoadString(params)
			if err != nil {
				return err
			}
		}
		s.OutputPlugin, err = creator(s.Mapping.Output.Config, s.Mapping.Output.Catalog, paramsMap)

		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("plugin %s not found", s.Mapping.Output.Plugin)
	}
	return nil
}
func (s *Session) InitInputPlugin() (err error) {
	if initializer, ok := s.InputPlugin.(plugins.Initializer); ok {
		err = initializer.Init()
		if err != nil {
			return err
		}
	}
	return nil
}
func (s *Session) InitOutputPlugin() (err error) {
	if initializer, ok := s.OutputPlugin.(plugins.Initializer); ok {
		err = initializer.Init()
		if err != nil {
			return err
		}
	}
	return nil
}
