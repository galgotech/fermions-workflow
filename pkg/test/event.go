package test

import (
	"context"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/serverlessworkflow/sdk-go/v2/model"

	"github.com/galgotech/fermions-workflow/pkg/worker/environment"
)

type eventStub struct {
	ch chan cloudevents.Event
}

func (e *eventStub) Produce(ctx context.Context, data model.Object) error {
	event := cloudevents.NewEvent()
	err := event.SetData(cloudevents.ApplicationJSON, model.ToInterface(data))
	if err != nil {
		return err
	}
	e.ch <- event
	return nil
}

func (e *eventStub) Consume(ctx context.Context) (cloudevents.Event, error) {
	event := <-e.ch
	return event, nil
}

var MapEvents = environment.MapEvents{
	"event0": &eventStub{ch: make(chan cloudevents.Event)},
}
