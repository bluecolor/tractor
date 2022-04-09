package runner

import (
	"context"
	"fmt"
	"net/rpc"
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
	isDriverDone    bool
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
	feedClient       *rpc.Client
}

func (r *Result) Errors() types.Errors {
	r.errors = types.Errors{}
	if r.isInputSuccess && r.isOutputSuccess && r.readCount != r.writeCount {
		r.errors.Add(fmt.Errorf("read count %d != write count %d", r.readCount, r.writeCount))
	}
	r.errors.Add(r.inputError)
	r.errors.Add(r.outputError)
	return r.errors
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
func (r *Result) IsSuccess() bool {
	if r.Errors().Count() > 0 {
		return false
	}
	if !r.isInputSuccess || !r.isOutputSuccess {
		return false
	}
	return true
}
func (r *Result) IsDone() bool {
	return r.isDone
}
func (r *Result) IsIODone() bool {
	return r.isInputDone && r.isOutputDone
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
		feedClient: GetFeedClient(options...),
		result: &Result{
			errors: types.Errors{},
		},
	}
	return r, nil
}

func (r *Runner) ProcessFeedback(f *msg.Feed) {
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
func (r *Runner) ProcessResult() *Result {
	log.Debug().Msgf("processing result ...")
	if !r.result.IsIODone() {
		log.Debug().Msgf("IO is not done yet")
		return r.result
	}
	if r.result.IsSuccess() {
		r.wire.SendSuccess(msg.Driver)
		r.wire.SendFeed(msg.NewSessionSuccess())
	}
	if r.result.Errors().Count() > 0 {
		r.wire.SendFeed(msg.NewSessionError())
		r.wire.SendFeed(msg.NewSessionDone())
	}
	return r.result
}
func (r *Runner) Run() (err error) {
	log.Debug().Msgf("runner started")
	r.wire.SendFeed(msg.NewSessionRunning())
	wg := &sync.WaitGroup{}
	wg.Add(3)
	// supervisor
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		err = r.Supervise(r.session.Extraction.GetTimeout()).Errors().Wrap()
	}(wg)
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
func (r *Runner) ForwardFeed(feed *msg.Feed) {
	if r.feedClient == nil {
		return
	}
	feed.SessionID = r.session.ID
	err := r.feedClient.Call("Handler.Process", feed, nil)
	if err != nil {
		log.Error().Msgf("feed client error %s", err)
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
			log.Debug().Msgf("received feedback %s", f)
			if !ok {
				err := fmt.Errorf("feedback channel closed unexpectedly")
				r.result.AddError(err, types.SupervisorError)
				result = r.Result()
				return
			}
			r.ForwardFeed(f)
			r.ProcessFeedback(f)
			r.TryCloseData()
			if r.IsIODone() {
				r.ProcessResult()
				if r.IsDone() {
					result = r.Result()
					r.TryCloseFeeds()
					return
				}
			}
		case <-r.ctx.Done():
			log.Debug().Msg("supervisor context done")
			r.Cancel()
		case <-time.After(timeout):
			log.Error().Msg("supervisor timeout")
			if !r.result.isTimeout {
				r.result.isTimeout = true
				r.TryCloseData()
			}
		}
	}
}
func (r *Runner) TryCloseFeeds() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.result.isInputDone && r.result.isOutputDone && !r.isFeedbackClosed {
		r.isFeedbackClosed = true
		r.wire.CloseFeeds()
		return true
	}
	return r.isFeedbackClosed
}
func (r *Runner) TryDone() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.result.isInputDone && r.result.isOutputDone && r.result.isDriverDone {
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
		r.isDataClosed = true
		return true
	}
	return r.isDataClosed
}
func (r *Runner) IsDone() bool {
	return r.TryDone()
}
func (r *Runner) IsIODone() bool {
	return r.result.IsIODone()
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
	r.wire.CloseFeeds()
}
