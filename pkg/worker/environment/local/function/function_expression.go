package function

import (
	"github.com/galgotech/fermions-workflow/pkg/worker/data"
	"github.com/galgotech/fermions-workflow/pkg/worker/filter"
)

func newExpression(operation string) *FunctionExpression {
	return &FunctionExpression{
		Operation: operation,
	}
}

type FunctionExpression struct {
	Operation string
	filter    filter.Filter
}

func (w *FunctionExpression) Init() (err error) {
	w.filter, err = filter.NewFilter(w.Operation)
	return err
}

func (w *FunctionExpression) Run(dataIn data.Data[any]) (dataOut data.Data[any], err error) {
	return w.filter.Run(dataIn)
}
