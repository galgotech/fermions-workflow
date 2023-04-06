package local

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/galgotech/fermions-workflow/pkg/test"
	"github.com/galgotech/fermions-workflow/pkg/worker/environment"
)

func TestLocal(t *testing.T) {
	local := &Local{}
	_ = environment.Environment(local)
}

func TestNext(t *testing.T) {
	t.Run("first state", func(t *testing.T) {
		local := &Local{}
		err := local.InitializeWorkflow(test.Workflow, test.NewBusStub())
		assert.Nil(t, err)

		// dataIn := data.Data[any]{}
		// dataOut, err := local.StateChannel(dataIn)
		// assert.NoError(t, err)
	})
}
