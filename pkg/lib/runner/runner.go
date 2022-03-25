package runner

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/bluecolor/tractor/pkg/lib/connectors"
	_ "github.com/bluecolor/tractor/pkg/lib/connectors/all"
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
	ctx              context.Context
	mu               sync.Mutex
	e                types.Extraction
	wire             *wire.Wire
	connectors       map[string]connectors.Connector
	result           *Result
	isFeedbackClosed bool
	isDataClosed     bool
	sessionID        string
	feedbackBackends []msg.FeedbackBackend
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

func New(ctx context.Context, e types.Extraction) (*Runner, error) {
	ic, err := connectors.GetConnector(
		e.SourceDataset.Connection.ConnectionType,
		connectors.ConnectorConfig(e.SourceDataset.Connection.Config),
	)
	if err != nil {
		return nil, err
	}
	inputConnector, ok := ic.(connectors.Input)
	if !ok {
		return nil, fmt.Errorf("connector %s is not an input connector", e.SourceDataset.Connection.ConnectionType)
	}

	oc, err := connectors.GetConnector(
		e.TargetDataset.Connection.ConnectionType,
		connectors.ConnectorConfig(e.TargetDataset.Connection.Config),
	)
	if err != nil {
		return nil, err
	}
	outputConnector, ok := oc.(connectors.Output)
	if !ok {
		return nil, fmt.Errorf("connector %s is not an output connector", e.TargetDataset.Connection.ConnectionType)
	}
	connectors := map[string]connectors.Connector{
		"input":  inputConnector,
		"output": outputConnector,
	}
	r := &Runner{
		ctx:              ctx,
		mu:               sync.Mutex{},
		wire:             wire.New(),
		e:                e,
		connectors:       connectors,
		feedbackBackends: make([]msg.FeedbackBackend, 0),
		result: &Result{
			errors: types.Errors{},
		},
	}
	return r, nil
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
func (r *Runner) Run(e types.Extraction, options ...Option) (err error) {
	r.SetOptions(options...)
	wg := &sync.WaitGroup{}
	wg.Add(3)
	// supervisor
	go func(wg *sync.WaitGroup, e types.Extraction) {
		defer wg.Done()
		err = r.Supervise(e.GetTimeout()).Eval().Errors().Wrap()
	}(wg, e)
	// output
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		r.RunOutput(*e.TargetDataset)
	}(wg)
	// input
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		r.RunInput(*e.SourceDataset)
	}(wg)
	wg.Wait()
	return
}
func (r *Runner) RunInput(d types.Dataset) error {
	defer func() {
		if err := r.inputConnector.Close(); err != nil {
			r.wire.SendInputError(err)
		}
	}()
	if err := r.inputConnector.Connect(); err != nil {
		r.wire.SendInputError(err)
		return err
	}
	return r.inputConnector.Read(d, r.wire)
}
func (r *Runner) RunOutput(d types.Dataset) error {

	if err := r.outputConnector.Connect(); err != nil {
		r.wire.SendOutputError(err)
		return err
	}

	return r.outputConnector.Write(d, r.wire)
}
func (r *Runner) ForwardFeedback(feedback *msg.Feedback) {
	for _, backend := range r.feedbackBackends {
		// ignore error
		backend.Store(r.sessionID, feedback)
	}
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
			r.ForwardFeedback(f)
			r.ProcessFeedback(f)
			r.TryCloseData()
			r.TryCloseFeedback()
			if r.IsDone() {
				result = r.Result()
				return
			}
		case <-r.ctx.Done():
			r.Cancel()
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
func (r *Runner) Cancel() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Result().AddError(fmt.Errorf("canceled"), types.SupervisorError)
	r.inputConnector.Close()
	r.outputConnector.Close()
	r.isDataClosed = true
	r.isFeedbackClosed = true
	r.wire.CloseData()
	r.wire.CloseFeedback()
}
