package function

import "github.com/galgotech/fermions-workflow/pkg/worker/data"

type FunctionGrpc struct {
	Operation string
}

func (w *FunctionGrpc) Init(data data.Data[any]) error {
	return nil
}

func (w *FunctionGrpc) Run(data data.Data[any]) (data.Data[any], error) {
	return nil, nil
}

func NewFunctionGrpc(operation string) *FunctionGrpc {
	return &FunctionGrpc{}
}
