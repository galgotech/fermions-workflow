package state

import (
	"context"
	"testing"

	"github.com/serverlessworkflow/sdk-go/v2/model"
	"github.com/stretchr/testify/assert"

	"github.com/galgotech/fermions-workflow/pkg/test"
)

func Test_newOperation(t *testing.T) {

	t.Run("prepare transition", func(t *testing.T) {
		stateBuilder := model.NewStateBuilder().Name("stateStart")
		stateBuilder.OperationState().
			ActionMode(model.ActionModeSequential).
			AddActions().Name("action1").FunctionRef().RefName("function1")
		state := stateBuilder.Build()

		baseState, err := NewBase(state, test.MapEvents)
		stateEnv, err := newOperation(*state.OperationState, baseState, test.MapFunctions, test.MapEvents)
		assert.NoError(t, err)
		assert.NotNil(t, state)
		assert.Equal(t, "stateStart", stateEnv.Name())
	})

	t.Run("run", func(t *testing.T) {
		stateBuilder := model.NewStateBuilder().Name("stateStart")
		stateBuilder.OperationState().
			ActionMode(model.ActionModeSequential).
			AddActions().Name("action1").FunctionRef().RefName("test0")
		state := stateBuilder.Build()

		baseState, err := NewBase(state, test.MapEvents)
		stateEnv, err := newOperation(*state.OperationState, baseState, test.MapFunctions, test.MapEvents)
		if assert.NoError(t, err) {
			dataIn := model.FromInterface(map[string]any{"test": "test"})
			dataOut, err := stateEnv.Run(context.Background(), dataIn)
			assert.NoError(t, err)
			assert.Equal(t, model.FromInterface(map[string]any{"test": "test"}), dataOut)
		}
	})
}
