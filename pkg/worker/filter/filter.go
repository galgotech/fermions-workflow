package filter

import (
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
	iter := f.query.Run(model.ToInterface(dataIn))
	v, ok := iter.Next()
	if !ok {
		return dataIn, nil
	}

	dataOut := model.FromInterface(v)
	return dataOut, nil
}

type EmptyFilter struct {
}

func (f *EmptyFilter) Run(dataIn model.Object) (model.Object, error) {
	return dataIn, nil
}
