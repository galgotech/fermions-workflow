package state

import (
	"context"
	"testing"

	"github.com/galgotech/fermions-workflow/pkg/test"
	"github.com/galgotech/fermions-workflow/pkg/worker/data"
	"github.com/galgotech/fermions-workflow/pkg/worker/environment"
	"github.com/stretchr/testify/assert"
)

func Test_newInject(t *testing.T) {

	t.Run("prepare state", func(t *testing.T) {
		mapEvents := environment.MapEvents{}
		baseState, err := NewBase(test.StateInject, mapEvents)
		assert.NoError(t, err)

		state, err := newInject(*test.StateInject.InjectState, baseState)
		assert.NoError(t, err)
		assert.NotNil(t, state)
		assert.Equal(t, "stateInject", state.Name())

		dataIn := data.Data[any]{}
		dataOut, err := state.Run(context.Background(), dataIn)
		assert.NoError(t, err)
		assert.Equal(t, dataIn, dataOut)
	})
}
