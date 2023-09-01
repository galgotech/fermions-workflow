package function

import (
	"github.com/galgotech/fermions-workflow/pkg/worker/filter"
	"github.com/serverlessworkflow/sdk-go/v2/model"
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

func (w *FunctionExpression) Run(dataIn model.Object) (dataOut model.Object, err error) {
	return w.filter.Run(dataIn)
}
