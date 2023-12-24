package state

import (
	"testing"

	"github.com/serverlessworkflow/sdk-go/v2/model"
	"github.com/stretchr/testify/assert"

	"github.com/galgotech/fermions-workflow/pkg/test"
)

func Test_newAction(t *testing.T) {
	builder := model.NewActionBuilder()
	builder.Name("action1")
	functionRef := builder.FunctionRef()
	functionRef.RefName("test0")

	actions, err := newAction([]model.Action{builder.Build()}, test.MapFunctions, test.MapEvents)
	assert.NoError(t, err)
	assert.NotNil(t, actions)
}

func TestRunAction(t *testing.T) {
	t.Run("useResults = true", func(t *testing.T) {
		actionBuilder1 := model.NewActionBuilder().Name("action1")
		actionBuilder1.FunctionRef().RefName("test0")
		actionBuilder2 := model.NewActionBuilder().Name("action2")
		actionBuilder2.FunctionRef().RefName("test1")

		actions, err := newAction([]model.Action{
			actionBuilder1.Build(),
			actionBuilder2.Build(),
		}, test.MapFunctions, test.MapEvents)
		assert.NoError(t, err)

		dataOut, err := actions.Run(model.FromMap(map[string]any{"input": "input_value"}))
		assert.NoError(t, err)
		assert.Equal(t, model.FromInterface(map[string]any{"input": "input_value", "out0": "value_out0", "out1": "value_out1"}), dataOut)
	})

	t.Run("useResults = false", func(t *testing.T) {
		actionBuilder1 := model.NewActionBuilder().Name("action1")
		actionBuilder1.FunctionRef().RefName("test0")
		actionBuilder2 := model.NewActionBuilder().Name("action2")
		actionBuilder2.FunctionRef().RefName("test1")
		actionBuilder2.ActionDataFilter().UseResults(true)

		actions, err := newAction([]model.Action{
			actionBuilder1.Build(),
			actionBuilder2.Build(),
		}, test.MapFunctions, test.MapEvents)
		if assert.NoError(t, err) {
			dataOut, err := actions.Run(model.FromMap(map[string]any{"input": "input_value"}))
			assert.NoError(t, err)
			assert.Equal(t, model.FromMap(map[string]any{"input": "input_value", "out2": "value_out2"}), dataOut)
		}
	})
}
