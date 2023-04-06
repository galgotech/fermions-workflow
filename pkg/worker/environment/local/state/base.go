package state

import (
	"context"
	"encoding/json"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"github.com/serverlessworkflow/sdk-go/v2/model"

	"github.com/galgotech/fermions-workflow/pkg/log"
	"github.com/galgotech/fermions-workflow/pkg/worker/data"
	"github.com/galgotech/fermions-workflow/pkg/worker/environment"
	"github.com/galgotech/fermions-workflow/pkg/worker/filter"
)

func NewBase(spec model.State, mapEvents environment.MapEvents) (s StateImpl, err error) {
	s.log = log.New("base state")

	produceEventSpec := []model.ProduceEvent{}
	if spec.Transition != nil {
		produceEventSpec = spec.Transition.ProduceEvents
	} else if spec.End != nil {
		produceEventSpec = spec.End.ProduceEvents
	}

	produceEvents := make([]*ProduceEventImp, len(produceEventSpec))
	for i, spec := range produceEventSpec {
		var f filter.Filter
		eventData := data.Data[any]{}
		if spec.Data.StrVal != "" {
			f, err = filter.NewFilter(spec.Data.StrVal)
			if err != nil {
				return s, err
			}
		} else {
			err := json.Unmarshal(spec.Data.RawValue, &eventData)
			if err != nil {
				return s, err
			}
		}

		produceEvents[i] = &ProduceEventImp{
			spec:   spec,
			event:  mapEvents[spec.EventRef],
			filter: f,
			data:   eventData,
		}
	}

	filterInput, filterOutput, err := initializeStateDataFilter(spec.StateDataFilter)
	if err != nil {
		return s, err
	}

	s.filterInput = filterInput
	s.filterOutput = filterOutput
	s.spec = spec
	s.produceEvents = produceEvents

	return s, nil
}

type ProduceEventImp struct {
	spec   model.ProduceEvent
	event  environment.Event
	filter filter.Filter
	data   data.Data[any]
}

func (p *ProduceEventImp) Data(dataIn data.Data[any]) (data.Data[any], error) {
	if p.filter != nil {
		return p.filter.Run(dataIn)
	}
	return p.data, nil
}

func (p *ProduceEventImp) Name() string {
	return p.spec.EventRef
}

type StateImpl struct {
	spec model.State
	log  log.Logger

	filterInput   filter.Filter
	filterOutput  filter.Filter
	produceEvents []*ProduceEventImp
}

func (p *StateImpl) Type() model.StateType {
	return p.spec.BaseState.Type
}

func (t *StateImpl) Name() string {
	return t.spec.Name
}

func (t *StateImpl) isEnd() bool {
	end := t.spec.End
	if end == nil {
		return false
	}
	return end.Terminate
}

func (t *StateImpl) Transition() string {
	return t.spec.Transition.NextState
}

func (t *StateImpl) FilterInput(data data.Data[any]) (data.Data[any], error) {
	return t.filterInput.Run(data)
}

func (t *StateImpl) FilterOutput(data data.Data[any]) (data.Data[any], error) {
	return t.filterOutput.Run(data)
}

func (t *StateImpl) ProduceEvents(ctx context.Context, dataIn data.Data[any]) error {
	for _, produceEvent := range t.produceEvents {
		dataOut, err := produceEvent.Data(dataIn)
		if err != nil {
			return err
		}

		// TODO: Create new event in data.Data[any]
		event := cloudevents.NewEvent()
		data, err := json.Marshal(dataOut)
		if err != nil {
			return err
		}
		err = event.UnmarshalJSON(data)
		if err != nil {
			return err
		}

		err = produceEvent.event.Produce(ctx, event)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *StateImpl) Next(ctx context.Context) <-chan environment.StateStart {
	ch := make(chan environment.StateStart)
	go func() {
		defer close(ch)
		if e.spec.End != nil && e.spec.End.Terminate {
			e.log.Info("state end", "name", e.spec.Name, "trace", ctx.Value("trace"))
			return
		}

		nextState := e.spec.Transition.NextState
		e.log.Debug("generate state next", "currentState", e.spec.Name, "nextState", nextState, "trace", ctx.Value("trace"))
		ch <- NewStateStartCtx(ctx, nextState)
	}()

	return ch
}

func NewStateStartCtx(ctx context.Context, state string) environment.StateStart {
	return &stateStart{ctx, state}
}

func NewStateStart(state string) (environment.StateStart, error) {
	trace, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "trace", trace.String())

	return &stateStart{ctx, state}, nil
}

type stateStart struct {
	ctx   context.Context
	state string
}

func (g *stateStart) Ctx() context.Context {
	return g.ctx
}

func (g *stateStart) State() string {
	return g.state
}
