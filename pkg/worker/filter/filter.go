package filter

import (
	"github.com/galgotech/fermions-workflow/pkg/worker/data"
	"github.com/itchyny/gojq"
	"github.com/serverlessworkflow/sdk-go/v2/model"
)

func NewFilter(filter string) (Filter, error) {
	if filter == "" {
		return &EmptyFilter{}, nil
	}

	query, err := gojq.Parse(filter)
	if err != nil {
		return nil, err
	}

	return &WorkflowFilter{
		query: query,
	}, nil
}

type Filter interface {
	Run(data model.Object) (model.Object, error)
}

type WorkflowFilter struct {
	query *gojq.Query
}

func (f *WorkflowFilter) Run(dataIn model.Object) (model.Object, error) {
	iter := f.query.Run(data.ToInterface(dataIn))
	v, ok := iter.Next()
	if !ok {
		return dataIn, nil
	}

	dataOut := data.FromInterface(v)
	return dataOut, nil
}

type EmptyFilter struct {
}

func (f *EmptyFilter) Run(dataIn model.Object) (model.Object, error) {
	return dataIn, nil
}
