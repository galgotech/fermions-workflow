package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProvideState(t *testing.T) {
	state := Provide()
	assert.NotNil(t, state)
}

func TestCurrentStateString(t *testing.T) {
	state := Provide()

	currentState := state.State("test")
	assert.NotNil(t, currentState)

	assert.Equal(t, nil, currentState)
}
