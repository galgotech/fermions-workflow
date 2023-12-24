package state

import (
	"context"
	"testing"

	"github.com/serverlessworkflow/sdk-go/v2/model"
	"github.com/stretchr/testify/assert"

	"github.com/galgotech/fermions-workflow/pkg/worker/environment"
)

func Test_newInject(t *testing.T) {

	t.Run("prepare state", func(t *testing.T) {
		stateBuilder := model.NewStateBuilder().Name("stateInject")
		stateBuilder.InjectState().Data(map[string]model.Object{
			"test0": model.FromString("testValStr"),
			"test1": model.FromInt(1),
			"test2": model.FromString("bytes"),
		})
		state := stateBuilder.Build()

		mapEvents := environment.MapEvents{}
		baseState, err := NewBase(state, mapEvents)
		assert.NoError(t, err)

		stateEnv, err := newInject(*state.InjectState, baseState)
		if assert.NoError(t, err) {
			assert.NotNil(t, state)
			assert.Equal(t, "stateInject", stateEnv.Name())
			dataOut, err := stateEnv.Run(context.Background(), model.Object{})
			if assert.NoError(t, err) {
				assert.Equal(t, model.Object{}, dataOut)
			}
		}

	})
}
