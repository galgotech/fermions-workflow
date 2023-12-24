package event

import (
	"context"
	"errors"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/serverlessworkflow/sdk-go/v2/model"

	"github.com/galgotech/fermions-workflow/pkg/bus"
	"github.com/galgotech/fermions-workflow/pkg/log"
	"github.com/galgotech/fermions-workflow/pkg/worker/environment"
)

func New(spec model.Event, busEvent bus.Bus) environment.Event {
	spec.DataOnly = true
	return &Event{
		log:      log.New("event"),
		spec:     spec,
		busEvent: busEvent,
	}
}

type Event struct {
	log      log.Logger
	spec     model.Event
	busEvent bus.Bus
}

func (e *Event) Consume(ctx context.Context) (cloudevents.Event, error) {
	if e.spec.Kind != model.EventKindConsumed {
		return cloudevents.NewEvent(), errors.New("event is not 'consume'")
	}

	e.log.Debug("consume", "source", e.spec.Source, "name", e.spec.Name, "type", e.spec.Type)
	event := <-e.busEvent.Subscribe(ctx, e.spec.Source)
	if event.Err != nil {
		return cloudevents.NewEvent(), event.Err
	}

	return event.Event, nil
}

func (e *Event) Produce(ctx context.Context, data model.Object) error {
	if e.spec.Kind != model.EventKindProduced {
		return errors.New("event is not 'produce'")
	}

	e.log.Debug("produce", "source", e.spec.Source, "name", e.spec.Name, "type", e.spec.Type)
	event := cloudevents.NewEvent()
	event.SetType(e.spec.Type)
	event.SetSource(e.spec.Source)
	err := event.SetData(cloudevents.ApplicationJSON, model.ToInterface(data))
	if err != nil {
		return err
	}

	e.busEvent.Publish(ctx, event)
	return nil
}
