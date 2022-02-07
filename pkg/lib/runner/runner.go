package runner

// import (
// 	"fmt"
// 	"sync"
// 	"time"

// 	"github.com/bluecolor/tractor/pkg/lib/connectors"
// 	_ "github.com/bluecolor/tractor/pkg/lib/connectors/all"
// 	"github.com/bluecolor/tractor/pkg/lib/meta"
// 	"github.com/bluecolor/tractor/pkg/lib/wire"
// 	"github.com/rs/zerolog/log"
// )

// type Runner struct {
// 	inputConnection  meta.Connection
// 	outputConnection meta.Connection
// 	wire             wire.Wire
// 	inputConnector   connectors.InputConnector
// 	outputConnector  connectors.OutputConnector
// 	Error            error
// }

// func New(inputConnection meta.Connection, outputConnection meta.Connection) (*Runner, error) {
// 	ic, err := connectors.GetConnector(
// 		outputConnection.ConnectionType,
// 		connectors.ConnectorConfig(inputConnection.Config),
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	inputConnector, ok := ic.(connectors.InputConnector)
// 	if !ok {
// 		return nil, fmt.Errorf("input connector %s is not an input connector", inputConnection.ConnectionType)
// 	}

// 	oc, err := connectors.GetConnector(
// 		outputConnection.ConnectionType,
// 		connectors.ConnectorConfig(outputConnection.Config),
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	outputConnector, ok := oc.(connectors.OutputConnector)
// 	if !ok {
// 		return nil, fmt.Errorf("output connector %s is not an output connector", outputConnection.ConnectionType)
// 	}

// 	return &Runner{
// 		inputConnection:  inputConnection,
// 		outputConnection: outputConnection,
// 		wire:             wire.New(),
// 		inputConnector:   inputConnector,
// 		outputConnector:  outputConnector,
// 		Error:            nil,
// 	}, nil
// }

// func (r *Runner) Run(p meta.ExtParams) error {
// 	wg := &sync.WaitGroup{}
// 	wg.Add(3)
// 	// supervisor
// 	go func(wg *sync.WaitGroup, p meta.ExtParams) {
// 		defer wg.Done()
// 		r.Supervise(p.GetTimeout())
// 	}(wg, p)
// 	// input
// 	go func(wg *sync.WaitGroup) {
// 		defer wg.Done()
// 		r.RunInput(p)
// 	}(wg)
// 	// output
// 	go func(wg *sync.WaitGroup) {
// 		defer wg.Done()
// 		r.RunOutput(p)
// 	}(wg)
// 	wg.Wait()
// 	return nil
// }

// func (r *Runner) RunInput(p meta.ExtParams) error {
// 	if err := r.inputConnector.Connect(); err != nil {
// 		r.wire.SendInputErrorFeed(err)
// 		return err
// 	}
// 	defer func() {
// 		if err := r.inputConnector.Close(); err != nil {
// 			r.wire.SendInputErrorFeed(err)
// 		}
// 	}()

// 	return r.inputConnector.Read(p, r.wire)
// }

// func (r *Runner) RunOutput(p meta.ExtParams) error {
// 	if err := r.outputConnector.Connect(); err != nil {
// 		r.wire.SendOutputErrorFeed(err)
// 		return err
// 	}
// 	defer func() {
// 		if err := r.inputConnector.Close(); err != nil {
// 			r.wire.SendInputErrorFeed(err)
// 		}
// 	}()

// 	return r.outputConnector.Write(p, r.wire)
// }

// func (r *Runner) Supervise(timeout time.Duration) error {
// 	success := 2
// 	for {
// 		select {
// 		case <-r.wire.IsReadDone():
// 			log.Info().Msg("read done")
// 			success--
// 			if success == 0 {
// 				r.wire.Done()
// 			}
// 		case <-r.wire.IsWriteDone():
// 			log.Info().Msg("write done")
// 			success--
// 			if success == 0 {
// 				r.wire.Done()
// 			}
// 		case <-r.wire.IsDone():
// 			log.Info().Msg("done")
// 			return nil
// 		case <-time.After(timeout):
// 			log.Info().Msg("timeout")
// 			err := fmt.Errorf("timeout")
// 			r.Error = err
// 			return err
// 		}
// 	}
// }
