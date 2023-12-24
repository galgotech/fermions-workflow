package state

import (
	"context"
	"sync"
	"testing"

	"github.com/serverlessworkflow/sdk-go/v2/builder"
	"github.com/serverlessworkflow/sdk-go/v2/model"
	"github.com/stretchr/testify/assert"

	"github.com/galgotech/fermions-workflow/pkg/test"
)

func Test_newEvent(t *testing.T) {
	stateBuilder := model.NewStateBuilder().
		Type(model.StateTypeEvent).
		Name("stateEvent0")
	stateBuilder.End().Terminate(true)
	eventStateBuilder := stateBuilder.EventState()
	onEvent := eventStateBuilder.AddOnEvents()
	onEvent.EventRefs([]string{"event0"})
	onEvent.AddActions().Name("action0").FunctionRef().RefName("test0")

	eventState := stateBuilder.Build()
	if !assert.NoError(t, builder.Validate(eventState)) {
		return
	}

	baseState, err := NewBase(eventState, test.MapEvents)
	if assert.NoError(t, err) {
		stateEnv, err := newEvent(eventStateBuilder.Build(), baseState, test.MapFunctions, test.MapEvents)
		if assert.NoError(t, err) {
			assert.NotNil(t, stateEnv)
			assert.Equal(t, "stateEvent0", stateEnv.Name())
		}
	}
}

func TestEvent(t *testing.T) {
	stateBuilder := model.NewStateBuilder().
		Type(model.StateTypeEvent).
		Name("stateEvent0")
	stateBuilder.End().Terminate(true)
	actionBuilder := stateBuilder.EventState().
		AddOnEvents().
		EventRefs([]string{"event0"}).
		AddActions().
		Name("action0")
	actionBuilder.FunctionRef().RefName("test0")

	eventState := stateBuilder.Build()
	if !assert.NoError(t, builder.Validate(eventState)) {
		return
	}

	for i := 0; i < 1; i++ {
		baseState, err := NewBase(eventState, test.MapEvents)
		assert.NoError(t, err)
		state, err := newEvent(*eventState.EventState, baseState, test.MapFunctions, test.MapEvents)
		assert.NoError(t, err)

		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			dataIn := model.FromMap(map[string]any{"input": "value_input"})
			dataOut, err := state.Run(context.Background(), dataIn)
			if assert.NoError(t, err) {
				dataExpected := model.FromMap(map[string]any{
					"input": "value_input",
					"out0":  "value_out0",
					"event": "value_event",
				})
				assert.Equal(t, dataExpected, dataOut)
			}
			wg.Done()
		}()

		err = test.MapEvents["event0"].Produce(context.Background(), model.FromMap(map[string]any{"event": "value_event"}))
		assert.NoError(t, err)

		wg.Wait()
	}
}
