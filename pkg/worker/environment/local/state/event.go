package state

import (
	"context"
	"errors"
	"fmt"

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
		actions, err := newAction(specEventRef.Actions, functions)
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
	eventsLen int
	onEvents  []EventRef
}

func (e *Event) Run(ctx context.Context, dataIn data.Data[any]) (data.Data[any], error) {
	eventOuts := make([]<-chan eventOut, len(e.onEvents))
	for i, onEvent := range e.onEvents {
		eventOut := e.runEvent(ctx, onEvent, dataIn)
		eventOuts[i] = eventOut
	}

	eventOut := <-concurrency.Or(eventOuts...)
	if eventOut.Err != nil {
		return nil, eventOut.Err
	}

	dataOut := eventOut.Data
	fmt.Println(dataOut)
	return dataOut, nil
}

func (e *Event) runEvent(ctx context.Context, onEvent EventRef, dataIn data.Data[any]) <-chan eventOut {
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
	Data    data.Data[any]
	Err     error
}

func NewStateEventGenerator(state environment.State) (environment.State, error) {
	stateEvent, ok := state.(*Event)
	if !ok {
		return nil, errors.New("state is not a Event")
	}

	return &stateEventGenerator{
		StateImpl: stateEvent.StateImpl,
		Event:     *stateEvent,
	}, nil
}

// When a workflow start is a StateEvent, keep waiting a event to start the workflow execution
type stateEventGenerator struct {
	StateImpl
	Event
}

func (s *stateEventGenerator) Next(ctx context.Context) <-chan environment.StateStart {
	ch := make(chan environment.StateStart)
	go func() {
		defer close(ch)
		ch <- <-s.Event.Next(ctx)

		stateStart, err := NewStateStart(s.spec.Name)
		if err != nil {
			s.log.Error("finished workflow not possive create the first state", "err", err.Error())
			return
		}
		s.log.Debug("generate state start", "startState", s.spec.Name, "trace", stateStart.Ctx().Value("trace"))
		ch <- stateStart
	}()

	return ch
}
