package state

import (
	"context"

	"github.com/serverlessworkflow/sdk-go/v2/model"

	"github.com/galgotech/fermions-workflow/pkg/concurrency"
	"github.com/galgotech/fermions-workflow/pkg/worker/data"
	"github.com/galgotech/fermions-workflow/pkg/worker/environment"
)

func newEvent(spec model.EventState, stateBase StateImpl, functions environment.MapFunctions, mapEvents environment.MapEvents) (environment.State, error) {
	onEvents := make([]EventRef, len(spec.OnEvents))
	for i, specEventRef := range spec.OnEvents {
		// Events
		events := make([]environment.Event, len(specEventRef.EventRefs))
		for j, eventRef := range specEventRef.EventRefs {
			events[j] = mapEvents[eventRef]
		}

		// Actions
		actions, err := newAction(specEventRef.Actions, functions, mapEvents)
		if err != nil {
			return nil, err
		}

		// Filters
		filterData, filterToStateData, err := initializeEventDataFilter(specEventRef.EventDataFilter)
		if err != nil {
			return nil, err
		}

		// Events
		e, err := newEventRef(events, actions, filterData, filterToStateData, spec.Exclusive)
		if err != nil {
			return nil, err
		}
		onEvents[i] = e
	}

	return &Event{
		StateImpl: stateBase,
		onEvents:  onEvents,
	}, nil
}

type Event struct {
	StateImpl
	onEvents []EventRef
}

func (e *Event) Run(ctx context.Context, dataIn model.Object) (model.Object, error) {
	eventOuts := make([]<-chan eventOut, len(e.onEvents))
	for i, onEvent := range e.onEvents {
		eventOut := e.runEvent(ctx, onEvent, dataIn)
		eventOuts[i] = eventOut
	}

	eventOut := <-concurrency.Or(eventOuts...)
	if eventOut.Err != nil {
		return data.ObjectNil, eventOut.Err
	}

	dataOut := eventOut.Data
	return dataOut, eventOut.Err
}

func (e *Event) runEvent(ctx context.Context, onEvent EventRef, dataIn model.Object) <-chan eventOut {
	ch := make(chan eventOut)
	go func() {
		defer close(ch)
		dataOut, err := onEvent.Run(ctx, dataIn)
		ch <- eventOut{
			OnEvent: onEvent,
			Data:    dataOut,
			Err:     err,
		}
	}()
	return ch
}

type eventOut struct {
	OnEvent EventRef
	Data    model.Object
	Err     error
}
