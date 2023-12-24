package state

import (
	"context"
	"sync"
	"testing"

	"github.com/serverlessworkflow/sdk-go/v2/builder"
	"github.com/serverlessworkflow/sdk-go/v2/model"
	"github.com/stretchr/testify/assert"

	"github.com/galgotech/fermions-workflow/pkg/test"
	"github.com/galgotech/fermions-workflow/pkg/worker/environment"
	"github.com/galgotech/fermions-workflow/pkg/worker/filter"
)

func TestRun(t *testing.T) {
	event := test.MapEvents["event0"]

	actionBuilder := model.NewActionBuilder().
		Name("action0")
	actionBuilder.FunctionRef().
		RefName("test0")
	action := actionBuilder.Build()
	err := builder.Validate(action)
	if !assert.NoError(t, err) {
		return
	}

	filterData, err := filter.NewFilter("")
	assert.NoError(t, err)

	events := []environment.Event{event}
	actions, err := newAction([]model.Action{action}, test.MapFunctions, test.MapEvents)
	assert.NoError(t, err)

	eventRef, err := newEventRef(events, actions, filterData, filterData, true)
	assert.NoError(t, err)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		dataInput := model.FromMap(map[string]any{})
		dataOut, err := eventRef.Run(context.Background(), dataInput)
		if assert.NoError(t, err) {
			expected := model.FromInterface(map[string]any{
				"out0":  "value_out0",
				"event": "value_event",
			})
			assert.Equal(t, expected, dataOut)
		}
		wg.Done()
	}()

	err = event.Produce(context.Background(), model.FromMap(map[string]any{"event": "value_event"}))
	if assert.NoError(t, err) {
		wg.Wait()
	}
}
