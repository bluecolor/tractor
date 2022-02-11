package runner

import (
	"fmt"
	"sync"
	"time"

	"github.com/bluecolor/tractor/pkg/lib/connectors"
	_ "github.com/bluecolor/tractor/pkg/lib/connectors/all"
	"github.com/bluecolor/tractor/pkg/lib/meta"
	"github.com/bluecolor/tractor/pkg/lib/msg"
	"github.com/bluecolor/tractor/pkg/lib/types"
	"github.com/bluecolor/tractor/pkg/lib/wire"
)

type Result struct {
	readCount       int
	writeCount      int
	isInputSuccess  bool
	isInputError    bool
	inputError      error
	isInputDone     bool
	isOutputSuccess bool
	isOutputError   bool
	outputError     error
	isOutputDone    bool
	isTimeout       bool
	isDone          bool
	errors          types.Errors
}

type Runner struct {
	mu               sync.Mutex
	inputConnection  meta.Connection
	outputConnection meta.Connection
	wire             *wire.Wire
	inputConnector   connectors.Input
	outputConnector  connectors.Output
	result           *Result
	isFeedbackClosed bool
	isDataClosed     bool
}

func (r *Result) Eval() *Result {
	// TODO: add other validations
	errs := types.Errors{}
	if r.isInputSuccess && r.isOutputSuccess && r.readCount != r.writeCount {
		errs.Add(fmt.Errorf("read count %d != write count %d", r.readCount, r.writeCount))
	}
	if errs.Count() > 0 {
		r.errors.Add(errs)
	}
	return r
}
func (r *Result) Errors() types.Errors {
	errs := types.Errors{}
	errs.Add(r.inputError)
	errs.Add(r.outputError)
	return errs
}

func (r *Result) AddError(err error, es ...types.ErrorSource) {
	if err == nil {
		return
	}
	source := types.UnknownErrorSource
	r.errors.Add(err)
	if len(es) != 0 {
		source = es[0]
	}
	switch source {
	case types.InputError:
		r.isInputError = true
		r.inputError = err
	case types.OutputError:
		r.isOutputError = true
		r.outputError = err
	}
	r.errors.Add(err)
}

func New(inputConnection meta.Connection, outputConnection meta.Connection) (*Runner, error) {
	ic, err := connectors.GetConnector(
		outputConnection.ConnectionType,
		connectors.ConnectorConfig(inputConnection.Config),
	)
	if err != nil {
		return nil, err
	}
	inputConnector, ok := ic.(connectors.Input)
	if !ok {
		return nil, fmt.Errorf("connector %s is not an input connector", inputConnection.ConnectionType)
	}

	oc, err := connectors.GetConnector(
		outputConnection.ConnectionType,
		connectors.ConnectorConfig(outputConnection.Config),
	)
	if err != nil {
		return nil, err
	}
	outputConnector, ok := oc.(connectors.Output)
	if !ok {
		return nil, fmt.Errorf("connector %s is not an output connector", outputConnection.ConnectionType)
	}

	return &Runner{
		mu:               sync.Mutex{},
		inputConnection:  inputConnection,
		outputConnection: outputConnection,
		wire:             wire.New(),
		inputConnector:   inputConnector,
		outputConnector:  outputConnector,
		result: &Result{
			errors: types.Errors{},
		},
	}, nil
}
func (r *Runner) ProcessFeedback(f *msg.Feedback) {
	r.result.readCount += f.InputProgress()
	r.result.writeCount += f.OutputProgress()
	switch {
	case f.IsInputSuccess():
		r.result.isInputSuccess = true
	case f.IsInputError(), f.IsOutputError():
		r.result.AddError(f.Error(), f.ErrorSource())
	case f.IsInputDone():
		r.result.isInputDone = true
	case f.IsOutputSuccess():
		r.result.isOutputSuccess = true
	case f.IsOutputDone():
		r.result.isOutputDone = true
	}
}
func (r *Runner) Result() *Result {
	return r.result
}
func (r *Runner) Run(p meta.ExtParams) (err error) {
	wg := &sync.WaitGroup{}
	wg.Add(3)
	// supervisor
	go func(wg *sync.WaitGroup, p meta.ExtParams) {
		defer wg.Done()
		err = r.Supervise(p.GetTimeout()).Eval().Errors().Wrap()
	}(wg, p)
	// output
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		r.RunOutput(p)
	}(wg)
	// input
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		r.RunInput(p)
	}(wg)
	wg.Wait()
	return
}
func (r *Runner) RunInput(p meta.ExtParams) error {
	defer func() {
		if err := r.inputConnector.Close(); err != nil {
			r.wire.SendInputError(err)
		}
	}()
	if err := r.inputConnector.Connect(); err != nil {
		r.wire.SendInputError(err)
		return err
	}
	return r.inputConnector.Read(p, r.wire)
}
func (r *Runner) RunOutput(p meta.ExtParams) error {

	if err := r.outputConnector.Connect(); err != nil {
		r.wire.SendOutputError(err)
		return err
	}

	return r.outputConnector.Write(p, r.wire)
}
func (r *Runner) Supervise(timeout time.Duration) (result *Result) {
	defer func() {
		if err := recover(); err != nil {
			r.result.AddError(fmt.Errorf("%v", err), types.SupervisorError)
		}
		result = r.result
	}()

	for {
		select {
		case f, ok := <-r.wire.ReceiveFeedback():
			if !ok {
				err := fmt.Errorf("feedback channel closed unexpectedly")
				r.result.AddError(err, types.SupervisorError)
				result = r.Result()
				return
			}
			fmt.Printf("feedback: %v\n", f)
			r.ProcessFeedback(f)
			r.TryCloseData()
			r.TryCloseFeedback()
			if r.IsDone() {
				result = r.Result()
				return
			}
		case <-time.After(timeout):
			if !r.result.isTimeout {
				r.result.isTimeout = true
				r.TryCloseData()
			}
		}
	}
}
func (r *Runner) TryCloseFeedback() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.result.isInputDone && r.result.isOutputDone && !r.isFeedbackClosed {
		r.isFeedbackClosed = true
		r.wire.CloseFeedback()
		return true
	}
	return r.isFeedbackClosed
}
func (r *Runner) TryDone() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.result.isInputDone && r.result.isOutputDone {
		r.result.isDone = true
		return true
	}
	return r.result.isDone
}
func (r *Runner) TryCloseData() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.result.isInputDone && !r.isDataClosed {
		r.wire.CloseData()
		return true
	}
	return r.isDataClosed
}
func (r *Runner) IsDone() bool {
	return r.TryDone()
}
