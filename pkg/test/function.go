package test

import (
	"github.com/serverlessworkflow/sdk-go/v2/model"

	"github.com/galgotech/fermions-workflow/pkg/worker/environment"
)

type FunctionStub struct {
	name  string
	value string
}

func (f *FunctionStub) Init() error {
	return nil
}

func (f *FunctionStub) Run(dataIn model.Object) (model.Object, error) {
	dataOut := model.FromInterface(map[string]any{f.name: f.value})
	return dataOut, nil
}

var Functions = []model.Function{
	{
		Name:      "test",
		Type:      "rest",
		Operation: "https://galgo.tech/functions.json#test",
	},
}

var MapFunctions = environment.MapFunctions{
	"test0": &FunctionStub{name: "out0", value: "value_out0"},
	"test1": &FunctionStub{name: "out1", value: "value_out1"},
}
