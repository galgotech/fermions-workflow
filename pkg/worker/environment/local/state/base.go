package state

import (
	"context"
	"encoding/json"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/serverlessworkflow/sdk-go/v2/model"

	"github.com/galgotech/fermions-workflow/pkg/log"
	"github.com/galgotech/fermions-workflow/pkg/worker/data"
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
		eventData := data.Data[any]{}
		if spec.Data.StrVal != "" {
			f, err = filter.NewFilter(spec.Data.StrVal)
			if err != nil {
				s.log.Error("fail load data 1", "data", spec.Data.StrVal)
				return s, err
			}
		} else {
			// err := json.Unmarshal(spec.Data.RawValue, &eventData)
			err := json.Unmarshal([]byte("{}"), &eventData)
			if err != nil {
				s.log.Error("fail load data 2", "data", spec.Data)
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

func (t *StateImpl) Transition() string {
	return t.spec.BaseState.Transition.NextState
}

func (t *StateImpl) FilterInput(data data.Data[any]) (data.Data[any], error) {
	return t.filterInput.Run(data)
}

func (t *StateImpl) FilterOutput(data data.Data[any]) (data.Data[any], error) {
	return t.filterOutput.Run(data)
}

func (t *StateImpl) ProduceEvents(ctx context.Context, dataIn data.Data[any]) error {
	for _, produceEvent := range t.produceEvents {
		// dataOut, err := produceEvent.Data(dataIn)
		// if err != nil {
		// 	return err
		// }

		// TODO: Create new event in data.Data[any]
		event := cloudevents.NewEvent()
		// data, err := json.Marshal(dataOut)
		// if err != nil {
		// 	t.log.Error("fail produce events", "dataOut", dataOut)
		// 	return err
		// }
		// err = event.UnmarshalJSON(data)
		// if err != nil {
		// 	return err
		// }

		err := produceEvent.event.Produce(ctx, event)
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
