package function

import (
	"testing"

	"github.com/galgotech/fermions-workflow/pkg/worker/data"
	"github.com/stretchr/testify/assert"
)

func TestExpression(t *testing.T) {
	e := newExpression(".")
	err := e.Init()
	assert.NoError(t, err)

	dataIn := data.Data[any]{"test": "test"}
	dataOut, err := e.Run(dataIn)
	assert.NoError(t, err)

	assert.Equal(t, dataIn, dataOut)
}
