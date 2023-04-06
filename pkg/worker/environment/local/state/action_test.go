package state

import (
	"testing"

	"github.com/galgotech/fermions-workflow/pkg/test"
	"github.com/galgotech/fermions-workflow/pkg/worker/data"
	"github.com/galgotech/fermions-workflow/pkg/worker/environment"
	"github.com/serverlessworkflow/sdk-go/v2/model"
	"github.com/stretchr/testify/assert"
)

type functionStub struct {
	environment.Function
}

func (f *functionStub) Run(dataIn data.Data[any]) (data.Data[any], error) {
	return data.Data[any]{"test": "test"}, nil
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

		dataOut, err := a.Run(data.Data[any]{})
		assert.Nil(t, err)
		assert.Equal(t, data.Data[any]{"test": "test"}, dataOut)
	})
}

func TestRunAction(t *testing.T) {
	t.Run("exec", func(t *testing.T) {
		mapFunctions := environment.MapFunctions{
			"test": &functionStub{},
		}

		a, err := newAction([]model.Action{test.Action1}, mapFunctions)
		assert.NoError(t, err)
		dataOut, err := a.Run(data.Data[any]{})
		assert.NoError(t, err)
		assert.Equal(t, data.Data[any]{"test": "test"}, dataOut)
	})
}
