package state

import (
	"context"
	"testing"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/serverlessworkflow/sdk-go/v2/model"
	"github.com/stretchr/testify/assert"

	"github.com/galgotech/fermions-workflow/pkg/test"
	"github.com/galgotech/fermions-workflow/pkg/worker/environment"
)

func Test_newEvent(t *testing.T) {
	mapFunctions := environment.MapFunctions{}
	mapFunctions["test"] = &functionStub{}
	mapEvents := environment.MapEvents{}
	mapEvents["event0"] = &eventStub{ch: make(chan cloudevents.Event)}

	baseState, _ := NewBase(test.EventState, mapEvents)
	state, err := newEvent(*test.EventState.EventState, baseState, mapFunctions, mapEvents)
	assert.Nil(t, err)
	assert.NotNil(t, state)
	assert.Equal(t, "event1", state.Name())

}

func TestEvent(t *testing.T) {
	mapFunctions := environment.MapFunctions{}
	mapFunctions["test"] = &functionStub{}
	mapEvents := environment.MapEvents{}
	mapEvents["nameEvent0"] = &eventStub{ch: make(chan cloudevents.Event)}

	baseState, _ := NewBase(test.EventState, mapEvents)
	state, _ := newEvent(*test.EventState.EventState, baseState, mapFunctions, mapEvents)

	go func() {
		event := cloudevents.NewEvent()
		event.SetSource("nameEvent0")
		event.SetType("type")
		event.SetData("application/json", map[string]string{"a": "b"})
		mapEvents["nameEvent0"].Produce(context.Background(), event)
	}()

	t.Run("run", func(t *testing.T) {
		dataIn := model.FromInterface(map[string]any{"test": "test"})
		dataOut, err := state.Run(context.Background(), dataIn)
		assert.NoError(t, err)
		assert.Equal(t, model.FromInterface(map[string]any{"test": "test"}), dataOut)
	})
}

type eventStub struct {
	environment.Event
	ch chan cloudevents.Event
}

func (e *eventStub) Produce(ctx context.Context, event cloudevents.Event) error {
	e.ch <- event
	return nil
}

func (e *eventStub) Consume(ctx context.Context) (cloudevents.Event, error) {
	event := <-e.ch
	return event, nil
}
