package data

import (
	"testing"

	"github.com/serverlessworkflow/sdk-go/v2/model"
	"github.com/stretchr/testify/assert"
)

func TestProvideState(t *testing.T) {
	state := Provide()
	assert.NotNil(t, state)
}

func TestState(t *testing.T) {
	state := Provide()

	currentState := state.Get()
	assert.Equal(t, model.FromMap(map[string]any{}), currentState)

	err := state.Set(model.FromMap(map[string]any{"test_key": "test_value"}))
	assert.NoError(t, err)
	assert.Equal(t, model.FromMap(map[string]any{"test_key": "test_value"}), state.Get())

	err = state.Set(model.FromMap(map[string]any{"test_key_2": 1}))
	assert.NoError(t, err)
	assert.Equal(t, model.FromMap(map[string]any{"test_key": "test_value", "test_key_2": 1}), state.Get())

	err = state.Set(model.FromString(""))
	assert.Error(t, err)
}
