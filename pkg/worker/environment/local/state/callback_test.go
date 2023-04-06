package state

import (
	"context"
	"testing"

	"github.com/galgotech/fermions-workflow/pkg/test"
	"github.com/galgotech/fermions-workflow/pkg/worker/data"
	"github.com/galgotech/fermions-workflow/pkg/worker/environment"
	"github.com/stretchr/testify/assert"
)

func Test_newCallback(t *testing.T) {

	t.Run("prepare state", func(t *testing.T) {
		mapFunctions := environment.MapFunctions{}
		mapFunctions["test"] = &functionStub{}
		mapEvents := environment.MapEvents{}

		baseState, err := NewBase(test.CallbackState, mapEvents)
		assert.NoError(t, err)

		state, err := newCallback(*test.CallbackState.CallbackState, baseState, mapFunctions, mapEvents)
		assert.Nil(t, err)
		assert.NotNil(t, state)
		assert.Equal(t, "stateCallback", state.Name())
	})
}

func TestCallback(t *testing.T) {
	mapFunctions := environment.MapFunctions{}
	mapFunctions["test"] = &functionStub{}
	mapEvents := environment.MapEvents{}

	baseState, err := NewBase(test.CallbackState, mapEvents)
	assert.NoError(t, err)

	state, _ := newCallback(*test.CallbackState.CallbackState, baseState, mapFunctions, mapEvents)

	t.Run("run", func(t *testing.T) {
		dataIn := data.Data[any]{"test": "test"}
		ctx := context.Background()

		dataOut, err := state.Run(ctx, dataIn)
		assert.NoError(t, err)
		assert.Equal(t, data.Data[any]{"test": "test"}, dataOut)
	})
}
