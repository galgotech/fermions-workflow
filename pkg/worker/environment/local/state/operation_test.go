package state

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/galgotech/fermions-workflow/pkg/test"
	"github.com/galgotech/fermions-workflow/pkg/worker/data"
	"github.com/galgotech/fermions-workflow/pkg/worker/environment"
)

func Test_newOperation(t *testing.T) {

	t.Run("prepare transition", func(t *testing.T) {
		mapFunctions := environment.MapFunctions{}
		mapFunctions["test"] = &functionStub{}
		mapEvents := environment.MapEvents{}

		baseState, err := NewBase(test.States[0], mapEvents)
		state, err := newOperation(*test.States[0].OperationState, baseState, mapFunctions)
		assert.NoError(t, err)
		assert.NotNil(t, state)
		assert.Equal(t, "stateStart", state.Name())
	})

	t.Run("run", func(t *testing.T) {
		mapFunctions := environment.MapFunctions{}
		mapFunctions["test"] = &functionStub{}
		mapEvents := environment.MapEvents{}

		baseState, err := NewBase(test.States[0], mapEvents)
		state, err := newOperation(*test.States[0].OperationState, baseState, mapFunctions)
		assert.NoError(t, err)

		dataIn := data.Data[any]{"test": "test"}
		dataOut, err := state.Run(context.Background(), dataIn)
		assert.NoError(t, err)
		assert.Equal(t, data.Data[any]{"test": "test"}, dataOut)
	})
}
