package function

import (
	"github.com/serverlessworkflow/sdk-go/v2/model"
)

type FunctionGrpc struct {
	Operation string
}

func (w *FunctionGrpc) Init(data model.Object) error {
	return nil
}

func (w *FunctionGrpc) Run(data model.Object) (model.Object, error) {
	return model.FromNull(), nil
}

func NewFunctionGrpc(operation string) *FunctionGrpc {
	return &FunctionGrpc{}
}
