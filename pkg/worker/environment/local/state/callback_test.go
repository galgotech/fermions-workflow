package state

import (
	"context"
	"testing"

	"github.com/serverlessworkflow/sdk-go/v2/model"
	"github.com/stretchr/testify/assert"

	"github.com/galgotech/fermions-workflow/pkg/test"
)

func Test_newCallback(t *testing.T) {

	t.Run("prepare state", func(t *testing.T) {
		stateBuilder := model.NewStateBuilder().Name("stateCallback")
		stateBuilder.CallbackState().
			Action().Name("action1").
			FunctionRef().RefName("test0")
		state := stateBuilder.Build()

		baseState, err := NewBase(state, test.MapEvents)
		assert.NoError(t, err)

		stateEnv, err := newCallback(*state.CallbackState, baseState, test.MapFunctions, test.MapEvents)
		if assert.NoError(t, err) {
			assert.NotNil(t, state)
			assert.Equal(t, "stateCallback", stateEnv.Name())
		}
	})
}

func TestCallback(t *testing.T) {
	stateBuilder := model.NewStateBuilder().Name("stateCallback")
	stateBuilder.CallbackState().
		Action().Name("action1").FunctionRef().RefName("test0")
	state := stateBuilder.Build()

	baseState, err := NewBase(state, test.MapEvents)
	assert.NoError(t, err)

	stateEnv, err := newCallback(*state.CallbackState, baseState, test.MapFunctions, test.MapEvents)
	if !assert.NoError(t, err) {
		return
	}

	t.Run("run", func(t *testing.T) {
		dataIn := model.FromInterface(map[string]any{"test": "test"})
		ctx := context.Background()

		dataOut, err := stateEnv.Run(ctx, dataIn)
		if assert.NoError(t, err) {
			assert.Equal(t, model.FromMap(map[string]any{"test": "test"}), dataOut)
		}
	})
}
