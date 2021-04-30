package cmd

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/bluecolor/tractor"
	cfg "github.com/bluecolor/tractor/config"
	"github.com/bluecolor/tractor/plugins/inputs"
	"github.com/bluecolor/tractor/plugins/outputs"
	"github.com/bluecolor/tractor/utils"
	"github.com/vbauerster/mpb/v6"
	"github.com/vbauerster/mpb/v6/decor"
)

type progress struct {
	show bool
	data *struct {
		wg      sync.WaitGroup
		total   int
		read    int
		written int
		rpb     *mpb.Bar
		wpb     *mpb.Bar
	}
}
type session struct {
	startTime time.Time
	mapping   *cfg.Mapping
	progress  *progress
	iplugin   tractor.Input
	oplugin   tractor.Output
	wire      tractor.Wire
}

func (s *session) start() (*sync.WaitGroup, error) {
	var wg sync.WaitGroup
	go func(wg *sync.WaitGroup) {
		s.iplugin.Read(s.wire)
		wg.Done()
		s.wire.CloseData()
	}(&wg)
	wg.Add(1)

	go func(wg *sync.WaitGroup) {
		s.oplugin.Write(s.wire)
		wg.Done()
		s.wire.CloseFeed()
	}(&wg)
	wg.Add(1)

	return &wg, nil
}

func (s *session) validateInput() error {
	if validator, ok := s.iplugin.(tractor.Validator); ok {
		if err := validator.Validate(); err != nil {
			return errors.New("❌  Failed to validate plugin config")
		} else {
			println("☑️  Plugin config validated")
			return nil
		}
	}
	return nil
}

func (s *session) validateOutput() error {
	if validator, ok := s.oplugin.(tractor.Validator); ok {
		if err := validator.Validate(); err != nil {
			return errors.New("❌  Failed to validate plugin config")
		} else {
			println("☑️  Plugin config validated")
			return nil
		}
	}
	return nil
}

func (s *session) createInputPlugin() (err error) {
	if creator, ok := inputs.Inputs[s.mapping.Input.Plugin]; ok {
		var params map[string]interface{} = nil
		if inputParams != "" {
			params, err = utils.JSONLoadString(inputParams)
			if err != nil {
				return err
			}
		}
		s.iplugin, err = creator(s.mapping.Input.Config, s.mapping.Input.Catalog, params)
		if err != nil {
			return err
		}
	} else {
		return errors.New(fmt.Sprintf("Plugin %s not found", s.mapping.Input.Plugin))
	}
	return nil
}

func (s *session) initInputPlugin() (err error) {
	if initializer, ok := s.iplugin.(tractor.Initializer); ok {
		err = initializer.Init()
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *session) createOutputPlugin(catalog *cfg.Catalog) (err error) {
	if creator, ok := outputs.Outputs[s.mapping.Output.Plugin]; ok {
		var params map[string]interface{} = nil
		if outputParams != "" {
			params, err = utils.JSONLoadString(outputParams)
			if err != nil {
				return err
			}
		}
		s.oplugin, err = creator(s.mapping.Output.Config, catalog, params)

		if err != nil {
			return err
		}
	} else {
		return errors.New(fmt.Sprintf("Plugin %s not found", s.mapping.Output.Plugin))
	}
	return nil
}

func (s *session) initOutputPlugin() (err error) {
	if initializer, ok := s.oplugin.(tractor.Initializer); ok {
		err = initializer.Init()
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *session) initMapping() (err error) {
	if mapping == "" {
		return errors.New("Mapping is not given")
	}
	conf := cfg.NewConfig()
	err = conf.LoadConfig(config)
	if err != nil {
		return err
	}
	s.mapping, err = conf.GetMapping(mapping)
	if err != nil {
		return err
	}
	return nil
}

func (s *session) initProgress() (err error) {
	s.progress = &progress{show: showProgress}

	if s.progress.show {
		s.progress.data = &struct {
			wg      sync.WaitGroup
			total   int
			read    int
			written int
			rpb     *mpb.Bar
			wpb     *mpb.Bar
		}{}
		if counter, ok := s.iplugin.(tractor.Counter); ok {
			s.progress.data.total, err = counter.Count()
			if err != nil {
				return err
			}
			p := mpb.New(mpb.WithWaitGroup(&s.progress.data.wg))
			s.progress.data.wg.Add(2)
			s.progress.data.rpb = p.AddBar(int64(s.progress.data.total),
				mpb.PrependDecorators(
					decor.Name("Read "),
					decor.Percentage(decor.WCSyncSpace),
				),
				mpb.AppendDecorators(
					decor.OnComplete(
						decor.EwmaETA(decor.ET_STYLE_GO, 60), "done",
					),
				),
			)
			s.progress.data.wpb = p.AddBar(int64(s.progress.data.total),
				mpb.PrependDecorators(
					decor.Name("Write"),
					decor.Percentage(decor.WCSyncSpace),
				),
				mpb.AppendDecorators(
					decor.OnComplete(
						decor.EwmaETA(decor.ET_STYLE_GO, 60), "done",
					),
				),
			)
		} else {
			s.progress.show = false
		}
	}
	return nil
}

func (s *session) init() (err error) {
	err = s.initMapping()
	if err != nil {
		return err
	}
	err = s.createInputPlugin()
	if err != nil {
		return err
	}
	err = s.initInputPlugin()
	if err != nil {
		return err
	}
	err = s.validateInput()
	if err != nil {
		return err
	}

	var catalog *cfg.Catalog = nil
	if s.mapping.Output.Catalog != nil {
		catalog = s.mapping.Output.Catalog
	} else {
		if d, ok := s.iplugin.(tractor.Discoverer); ok {
			catalog, err = d.Discover()
			if err != nil {
				return err
			}
		}
	}
	err = s.createOutputPlugin(catalog)
	if err != nil {
		return err
	}

	err = s.initOutputPlugin()
	if err != nil {
		return err
	}
	err = s.validateOutput()
	if err != nil {
		return err
	}

	err = s.initProgress()
	if err != nil {
		return err
	}

	return nil
}

func (s *session) end() {
	if s.progress.show {
		if s.progress.data.rpb != nil {
			s.progress.data.rpb.Abort(true)
		}
		if s.progress.data.wpb != nil {
			s.progress.data.wpb.Abort(true)
		}
	}
}

func (s *session) listenFeeds() {
	for f := range s.wire.ReadFeeds() {
		switch f.Type {
		case tractor.Progress:
			if s.progress.show {
				s.processProgressFeed(f)
			}
		case tractor.Success:
			if s.progress.show {
				s.progress.data.wg.Done()
			}
		case tractor.Error:
			if s.progress.show {
				s.progress.data.wg.Done()
			}
		}
	}
}

func (s *session) processProgressFeed(f tractor.Feed) {
	p := f.Content.(tractor.ProgressFeed)
	switch f.Sender {
	case tractor.InputPlugin:
		s.progress.data.read += p.Count()
		s.progress.data.rpb.IncrBy(p.Count())
	case tractor.OutputPlugin:
		s.progress.data.written += p.Count()
		s.progress.data.wpb.IncrBy(p.Count())
	}
}

func (s *session) duration() time.Duration {
	return time.Since(s.startTime)
}
