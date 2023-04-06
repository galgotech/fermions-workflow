package filter

import (
	"errors"

	"github.com/itchyny/gojq"

	"github.com/galgotech/fermions-workflow/pkg/worker/data"
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
	Run(data data.Data[any]) (data.Data[any], error)
}

type WorkflowFilter struct {
	query *gojq.Query
}

func (f *WorkflowFilter) Run(dataIn data.Data[any]) (dataOut data.Data[any], err error) {
	iter := f.query.Run(dataIn.ToMap())
	v, ok := iter.Next()
	if !ok {
		return dataIn, nil
	}
	if val, ok := v.(map[string]any); ok {
		dataOut.FromMap(val)
	} else {
		return nil, errors.New("invalid filter response")
	}

	return dataOut, err
}

type EmptyFilter struct {
}

func (f *EmptyFilter) Run(dataIn data.Data[any]) (data.Data[any], error) {
	return dataIn, nil
}
