package session

import (
	"errors"
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

	_ "github.com/bluecolor/tractor/lib/plugins/inputs/all"
	_ "github.com/bluecolor/tractor/lib/plugins/outputs/all"
)

type Progress struct {
	Show bool
	Data *struct {
		Total   int
		Read    int
		Written int
	}
}

type Session struct {
	StartTime    time.Time
	EndTime      time.Time
	Mapping      *config.Mapping
	Progress     *Progress
	InputPlugin  inputs.InputPlugin
	OutputPlugin outputs.OutputPlugin
	Wire         *wire.Wire
}

func (s *Session) Duration() time.Duration {
	return s.EndTime.Sub(s.StartTime)
}
func (s *Session) createInputPlugin(params string) (err error) {
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
func (s *Session) createOutputPlugin(catalog *config.Catalog, params string) (err error) {
	if creator, ok := outputs.Outputs[s.Mapping.Output.Plugin]; ok {
		var paramsMap map[string]interface{} = nil
		if params != "" {
			paramsMap, err = utils.JSONLoadString(params)
			if err != nil {
				return err
			}
		}
		s.OutputPlugin, err = creator(s.Mapping.Output.Config, catalog, paramsMap)

		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("plugin %s not found", s.Mapping.Output.Plugin)
	}
	return nil
}
func (s *Session) initInputPlugin() (err error) {
	if initializer, ok := s.InputPlugin.(plugins.Initializer); ok {
		err = initializer.Init()
		if err != nil {
			return err
		}
	}
	return nil
}
func (s *Session) initOutputPlugin() (err error) {
	if initializer, ok := s.OutputPlugin.(plugins.Initializer); ok {
		err = initializer.Init()
		if err != nil {
			return err
		}
	}
	return nil
}
func (s *Session) initMapping(conf *config.Config, mapping string) (err error) {
	if mapping == "" {
		return errors.New("mapping is not given")
	}
	if err != nil {
		return err
	}
	s.Mapping, err = conf.GetMapping(mapping)
	if err != nil {
		return err
	}
	return nil
}
func (s *Session) initSessionProgress(showProgress bool) (err error) {
	s.Progress = &Progress{Show: showProgress}
	return
}
func (s *Session) initWire() {
	s.Wire = wire.NewWire()
}
func (s *Session) processProgressFeed(f feed.Feed) {
	p := f.Content.(feed.Progress)
	switch f.Sender {
	case feed.SenderInputPlugin:
		s.Progress.Data.Read += p.Count()
	case feed.SenderOutputPlugin:
		s.Progress.Data.Written += p.Count()
	}
}
func (s *Session) listen() {
	for f := range s.Wire.ReadFeed() {
		switch f.Type {
		case feed.ProgressFeed:
			if s.Progress.Show {
				s.processProgressFeed(f)
			}
		case feed.SuccessFeed:
			fmt.Println("Success") // todo
		case feed.ErrorFeed:
			fmt.Println("Error") // todo
		}
	}
}
func (s *Session) prepare(conf *config.Config, mapping string, showProgress bool, params string) (err error) {
	s.initWire()
	err = s.initMapping(conf, mapping)
	if err != nil {
		return err
	}
	err = s.createInputPlugin(params)
	if err != nil {
		return err
	}
	err = s.initInputPlugin()
	if err != nil {
		return err
	}
	err = s.validateInputPlugin()
	if err != nil {
		return err
	}

	var catalog *config.Catalog = nil
	if s.Mapping.Output.Catalog != nil {
		catalog = s.Mapping.Output.Catalog
	} else {
		if d, ok := s.InputPlugin.(plugins.CatalogDiscoverer); ok {
			catalog, err = d.DiscoverCatalog()
			if err != nil {
				return err
			}
		}
	}
	err = s.createOutputPlugin(catalog, params)
	if err != nil {
		return err
	}

	err = s.initOutputPlugin()
	if err != nil {
		return err
	}
	err = s.validateOutputPlugin()
	if err != nil {
		return err
	}

	err = s.initSessionProgress(showProgress)
	if err != nil {
		return err
	}

	return nil
}
func (s *Session) start() (*sync.WaitGroup, error) {
	s.StartTime = time.Now()
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
	}(&wg)
	wg.Add(1)

	return &wg, nil
}
func (s *Session) end() {
	s.EndTime = time.Now()
	s.Wire.CloseFeed()
}

func (s *Session) Run() (err error) {
	wg, err := s.start()
	if err != nil {
		return
	}
	go s.listen()
	wg.Wait()
	s.end()
	return
}
func NewSession(conf *config.Config, mapping string, showProgress bool, params string) (s *Session, err error) {
	s = &Session{}
	err = s.prepare(conf, mapping, showProgress, params)
	if err != nil {
		return
	}
	return s, nil
}
