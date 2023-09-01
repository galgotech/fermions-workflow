package state

import (
	"testing"

	"github.com/galgotech/fermions-workflow/pkg/test"
	"github.com/galgotech/fermions-workflow/pkg/worker/environment"
	"github.com/serverlessworkflow/sdk-go/v2/model"
	"github.com/stretchr/testify/assert"
)

type functionStub struct {
	environment.Function
}

func (f *functionStub) Run(dataIn model.Object) (model.Object, error) {
	dataOut := model.FromInterface(map[string]any{"test": "test"})
	return dataOut, nil
}

func Test_newAction(t *testing.T) {

	t.Run("prepare transition", func(t *testing.T) {
		mapFunctions := environment.MapFunctions{
			"test": &functionStub{},
		}

		a, err := newAction([]model.Action{test.Action1}, mapFunctions)
		assert.NoError(t, err)
		assert.NotNil(t, a)
	})

	t.Run("run", func(t *testing.T) {
		mapFunctions := environment.MapFunctions{
			"test": &functionStub{},
		}
		a, err := newAction([]model.Action{test.Action1}, mapFunctions)
		assert.NoError(t, err)

		dataOut, err := a.Run(model.Object{})
		assert.Nil(t, err)
		assert.Equal(t, model.FromInterface(map[string]any{"test": "test"}), dataOut)
	})
}

func TestRunAction(t *testing.T) {
	t.Run("exec", func(t *testing.T) {
		mapFunctions := environment.MapFunctions{
			"test": &functionStub{},
		}

		a, err := newAction([]model.Action{test.Action1}, mapFunctions)
		assert.NoError(t, err)
		dataOut, err := a.Run(model.Object{})
		assert.NoError(t, err)
		assert.Equal(t, model.FromInterface(map[string]any{"test": "test"}), dataOut)
	})
}
