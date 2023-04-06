package state

import (
	"context"

	"github.com/galgotech/fermions-workflow/pkg/concurrency"
	"github.com/galgotech/fermions-workflow/pkg/worker/data"
	"github.com/galgotech/fermions-workflow/pkg/worker/environment"
	"github.com/galgotech/fermions-workflow/pkg/worker/filter"
)

func newEventRef(events []environment.Event, actions Actions, filterData filter.Filter, filterToStateData filter.Filter, exclusive bool) (e EventRef, err error) {
	e.filterData = filterData
	e.filterToStateData = filterToStateData
	e.events = events
	e.actions = actions
	e.exclusive = exclusive
	return e, nil
}

type EventRef struct {
	filterData        filter.Filter
	filterToStateData filter.Filter

	exclusive bool
	events    []environment.Event
	actions   Actions
}

func (e *EventRef) Run(ctx context.Context, dataIn data.Data[any]) (data.Data[any], error) {
	dataIn, err := e.filterData.Run(dataIn)
	if err != nil {
		return nil, err
	}

	dataIn, err = e.filterToStateData.Run(dataIn)
	if err != nil {
		return nil, err
	}

	eventRefs := make([]<-chan eventRefOut, len(e.events))
	for i, event := range e.events {
		eventRefs[i] = e.runEventRef(ctx, event)
	}

	if e.exclusive {
		eventRef := <-concurrency.Or(eventRefs...)
		if eventRef.Err != nil {
			return nil, eventRef.Err
		}
		dataIn.Merge(eventRef.Data)

	} else {
		for _, ch := range eventRefs {
			eventRef := <-ch
			if eventRef.Err != nil {
				return nil, eventRef.Err
			}
			dataIn.Merge(eventRef.Data)
		}
	}

	dataOut, err := e.actions.Run(dataIn)
	if err != nil {
		return nil, err
	}

	return dataOut, err
}

func (e *EventRef) runEventRef(ctx context.Context, event environment.Event) <-chan eventRefOut {
	ch := make(chan eventRefOut)
	go func() {
		defer close(ch)
		event, err := event.Consume(ctx)
		if err != nil {
			ch <- eventRefOut{
				Data: nil,
				Err:  err,
			}
			return
		}

		dataOut, err := data.FromEvent(event)
		ch <- eventRefOut{
			Data: dataOut,
			Err:  err,
		}
	}()

	return ch
}

type eventRefOut struct {
	Data data.Data[any]
	Err  error
}
