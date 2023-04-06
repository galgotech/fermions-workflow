package state

import (
	"context"

	"github.com/serverlessworkflow/sdk-go/v2/model"

	"github.com/galgotech/fermions-workflow/pkg/worker/data"
	"github.com/galgotech/fermions-workflow/pkg/worker/environment"
)

func newCallback(spec model.CallbackState, baseState StateImpl, functions environment.MapFunctions, mapEvents environment.MapEvents) (environment.State, error) {
	// Events
	event := []environment.Event{mapEvents[spec.EventRef]}

	// Actions
	actions, err := newAction([]model.Action{spec.Action}, functions)
	if err != nil {
		return nil, err
	}

	// Filters
	filterData, filterToStateData, err := initializeEventDataFilter(*spec.EventDataFilter)
	if err != nil {
		return nil, err
	}

	eventRef, err := newEventRef(event, actions, filterData, filterToStateData, true)
	if err != nil {
		return nil, err
	}

	c := &Callback{
		StateImpl: baseState,
		EventRef:  eventRef,
	}

	return c, nil
}

type Callback struct {
	StateImpl
	EventRef
}

func (c *Callback) Run(ctx context.Context, dataIn data.Data[any]) (dataOut data.Data[any], err error) {
	return dataIn, nil
}
