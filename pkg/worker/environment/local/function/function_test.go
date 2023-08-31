package function

import (
	"testing"

	"github.com/serverlessworkflow/sdk-go/v2/model"
	"github.com/stretchr/testify/assert"
)

func TestErrFunctionTypeInvalid(t *testing.T) {
	_, err := New(model.Function{
		Type: "invalid",
	})
	assert.Error(t, err)
}

func TestPrepareFunctionRest(t *testing.T) {
	functionRest, err := New(model.Function{
		Type: model.FunctionTypeREST,
	})

	assert.NotNil(t, functionRest)
	assert.Nil(t, err)
}
