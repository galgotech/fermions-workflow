package state

import (
	"context"
	"errors"

	"github.com/serverlessworkflow/sdk-go/v2/model"

	"github.com/galgotech/fermions-workflow/pkg/log"
	"github.com/galgotech/fermions-workflow/pkg/worker/environment"
	"github.com/galgotech/fermions-workflow/pkg/worker/filter"
)

func NewBase(spec model.State, mapEvents environment.MapEvents) (s StateImpl, err error) {
	s.log = log.New("base state")

	produceEventSpec := []model.ProduceEvent{}
	if spec.BaseState.Transition != nil {
		produceEventSpec = spec.BaseState.Transition.ProduceEvents
	} else if spec.BaseState.End != nil {
		produceEventSpec = spec.BaseState.End.ProduceEvents
	}

	produceEvents := make([]*ProduceEventImp, len(produceEventSpec))
	for i, spec := range produceEventSpec {
		var f filter.Filter
		eventData := model.FromNull()
		if spec.Data.Type == model.Map {
			eventData = spec.Data
		} else if spec.Data.Type == model.String {
			f, err = filter.NewFilter(spec.Data.StringValue)
			if err != nil {
				s.log.Error("fail load data 1", "data", spec.Data.StringValue)
				return s, err
			}
		} else {
			return s, errors.New("invalid produceEvent.data")
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
	data   model.Object
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

func (t *StateImpl) Transition() string {
	return t.spec.BaseState.Transition.NextState
}

func (t *StateImpl) FilterInput(data model.Object) (model.Object, error) {
	return t.filterInput.Run(data)
}

func (t *StateImpl) FilterOutput(data model.Object) (model.Object, error) {
	return t.filterOutput.Run(data)
}

func (t *StateImpl) ProduceEvents(ctx context.Context, dataIn model.Object) (err error) {
	for _, produceEvent := range t.produceEvents {
		if produceEvent.data.Type == model.Null {
			dataIn, err = produceEvent.filter.Run(dataIn)
			if err != nil {
				return err
			}
		} else {
			dataIn = produceEvent.data
		}

		err = produceEvent.event.Produce(ctx, dataIn)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *StateImpl) Next() (string, bool) {
	if e.spec.End != nil && e.spec.End.Terminate {
		return "", false
	}

	return e.spec.BaseState.Transition.NextState, true
}
