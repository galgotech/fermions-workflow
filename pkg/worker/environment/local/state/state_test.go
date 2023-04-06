package state

import (
	"testing"

	"github.com/galgotech/fermions-workflow/pkg/worker/environment"
	"github.com/serverlessworkflow/sdk-go/v2/model"
	"github.com/stretchr/testify/assert"
)

func TestPrepareState(t *testing.T) {
	states := []model.State{
		{
			BaseState: model.BaseState{
				Type: model.StateTypeOperation,
				Name: "stateStart",
				End: &model.End{
					Terminate: true,
				},
			},
			OperationState: &model.OperationState{
				ActionMode: model.ActionModeSequential,
				Actions: []model.Action{{
					Name: "action1",
				}},
			},
		},
	}

	t.Run("prepare state", func(t *testing.T) {
		transitionPrepared, err := New(states[0], make(environment.MapFunctions, 0), make(environment.MapEvents, 0))
		assert.Nil(t, err)
		assert.NotNil(t, transitionPrepared)
		assert.IsType(t, &Operation{}, transitionPrepared)
	})
}
