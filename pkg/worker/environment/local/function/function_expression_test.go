package function

import (
	"testing"

	"github.com/serverlessworkflow/sdk-go/v2/model"
	"github.com/stretchr/testify/assert"
)

func TestExpression(t *testing.T) {
	e := newExpression(".")
	err := e.Init()
	assert.NoError(t, err)

	dataIn := model.FromInterface(any(map[string]any{"test": "test"}))
	dataOut, err := e.Run(dataIn)
	assert.NoError(t, err)
	assert.Equal(t, dataIn, dataOut)

	e = newExpression(".test")
	err = e.Init()
	assert.NoError(t, err)

	dataIn = model.FromInterface(any(map[string]any{"test": "test"}))
	dataOut, err = e.Run(dataIn)
	assert.NoError(t, err)
	assert.Equal(t, model.FromInterface("test"), dataOut)
}
