package event

import (
	"context"
	"encoding/json"
	"sync"
	"testing"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/serverlessworkflow/sdk-go/v2/model"
	"github.com/stretchr/testify/assert"

	"github.com/galgotech/fermions-workflow/pkg/test"
)

var EventProduced = model.Event{
	Kind:     model.EventKindProduced,
	Name:     "EventName",
	Type:     "eventType",
	Source:   "eventSource",
	DataOnly: true,
}

var EventConsumed = model.Event{
	Kind:     model.EventKindConsumed,
	Name:     "EventName",
	Type:     "eventType",
	Source:   "eventSource",
	DataOnly: true,
}

func TestProduce(t *testing.T) {

	busStub := test.NewBusStub()
	event := New(EventProduced, busStub)

	var wg sync.WaitGroup
	wg.Add(1)

	ctx := context.Background()
	ch := busStub.Subscribe(ctx, "eventSource")
	go func() {
		select {
		case <-time.After(1 * time.Second):
			assert.Fail(t, "timeout")
		case v := <-ch:
			if assert.NoError(t, v.Err) {
				var object any
				err := json.Unmarshal(v.Event.Data(), &object)
				assert.NoError(t, err)
				assert.Equal(t, map[string]any{"test-key": "test-value"}, object)
			}
		}
		wg.Done()
	}()

	event.Produce(ctx, model.FromMap(map[string]any{"test-key": "test-value"}))

	wg.Wait()
}

func TestConsume(t *testing.T) {
	busStub := test.NewBusStub()
	workflowEvent := New(EventConsumed, busStub)

	var wg sync.WaitGroup
	wg.Add(1)

	ctx := context.Background()

	go func() {
		event, err := workflowEvent.Consume(ctx)
		if assert.NoError(t, err) {
			var object any
			err := json.Unmarshal(event.Data(), &object)
			assert.NoError(t, err)
			assert.Equal(t, map[string]any{"test-key": "test-value"}, object)
		}
		wg.Done()
	}()
	time.Sleep(200 * time.Millisecond)

	event := cloudevents.NewEvent()
	event.SetType("eventType")
	event.SetSource("eventSource")
	err := event.SetData(cloudevents.ApplicationJSON, map[string]any{"test-key": "test-value"})
	assert.NoError(t, err)

	busStub.Publish(ctx, event)

	wg.Wait()
}
