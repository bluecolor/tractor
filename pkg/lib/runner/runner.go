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
	"github.com/rs/zerolog/log"
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

type ioConnectors struct {
	input  connectors.Input
	output connectors.Output
}

type Runner struct {
	ctx              context.Context
	mu               sync.Mutex
	session          types.Session
	wire             *wire.Wire
	connectors       ioConnectors
	result           *Result
	isFeedbackClosed bool
	isDataClosed     bool
	feedBackends     []msg.FeedBackend
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

func New(ctx context.Context, s types.Session, options ...Option) (*Runner, error) {
	ic, err := connectors.GetConnector(
		s.Extraction.SourceDataset.Connection.ConnectionType,
		connectors.ConnectorConfig(s.Extraction.SourceDataset.Connection.Config),
	)
	if err != nil {
		return nil, err
	}
	inputConnector, ok := ic.(connectors.Input)
	if !ok {
		return nil, fmt.Errorf(
			"connector %s is not an input connector",
			s.Extraction.SourceDataset.Connection.ConnectionType,
		)
	}

	oc, err := connectors.GetConnector(
		s.Extraction.TargetDataset.Connection.ConnectionType,
		connectors.ConnectorConfig(s.Extraction.TargetDataset.Connection.Config),
	)
	if err != nil {
		return nil, err
	}
	outputConnector, ok := oc.(connectors.Output)
	if !ok {
		return nil, fmt.Errorf(
			"connector %s is not an output connector",
			s.Extraction.TargetDataset.Connection.ConnectionType,
		)
	}
	r := &Runner{
		ctx:     ctx,
		mu:      sync.Mutex{},
		wire:    wire.New(),
		session: s,
		connectors: ioConnectors{
			input:  inputConnector,
			output: outputConnector,
		},
		feedBackends: GetFeedBackends(options...),
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
		log.Info().Msgf("input success")
		r.result.isInputSuccess = true
	case f.IsInputError(), f.IsOutputError():
		log.Info().Msgf("input error %s", f.Error())
		r.result.AddError(f.Error(), f.ErrorSource())
	case f.IsInputDone():
		log.Info().Msgf("input done")
		r.result.isInputDone = true
	case f.IsOutputSuccess():
		log.Info().Msgf("output success")
		r.result.isOutputSuccess = true
	case f.IsOutputDone():
		log.Info().Msgf("output done")
		r.result.isOutputDone = true
	}
}
func (r *Runner) Result() *Result {
	return r.result
}
func (r *Runner) Run() (err error) {
	log.Info().Msgf("runner started")
	r.wire.SendFeedback(msg.NewRunning())
	wg := &sync.WaitGroup{}
	wg.Add(3)
	// supervisor
	go func(wg *sync.WaitGroup, s types.Session) {
		defer wg.Done()
		err = r.Supervise(s.Extraction.GetTimeout()).Eval().Errors().Wrap()
	}(wg, r.session)
	// output
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		r.RunOutput(*r.session.Extraction.TargetDataset)
	}(wg)
	// input
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		r.RunInput(*r.session.Extraction.SourceDataset)
	}(wg)
	wg.Wait()
	return
}
func (r *Runner) RunInput(d types.Dataset) error {
	defer func() {
		if err := r.connectors.input.Close(); err != nil {
			r.wire.SendInputError(err)
		}
	}()
	if err := r.connectors.input.Connect(); err != nil {
		log.Error().Msgf("input connect error %s", err)
		r.wire.SendInputError(err)
		return err
	}
	return r.connectors.input.Read(d, r.wire)
}
func (r *Runner) RunOutput(d types.Dataset) error {
	defer func() {
		if err := r.connectors.output.Close(); err != nil {
			r.wire.SendOutputError(err)
		}
	}()
	if err := r.connectors.output.Connect(); err != nil {
		r.wire.SendOutputError(err)
		return err
	}

	return r.connectors.output.Write(d, r.wire)
}
func (r *Runner) ForwardFeedback(feedback *msg.Feedback) {
	for _, backend := range r.feedBackends {
		// ignore error
		backend.Store(r.session.ID, feedback)
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
	r.connectors.input.Close()
	r.connectors.output.Close()
	r.isDataClosed = true
	r.isFeedbackClosed = true
	r.wire.CloseData()
	r.wire.CloseFeedback()
}
