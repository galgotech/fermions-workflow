package state

import (
	"context"

	"github.com/galgotech/fermions-workflow/pkg/worker/data"
	"github.com/serverlessworkflow/sdk-go/v2/model"
)

func newInject(spec model.InjectState, baseState StateImpl) (*Inject, error) {
	dataIn := data.Data[any]{}
	for key, dataOut := range spec.Data {
		switch dataOut.Type {
		case model.String:
			dataIn[key] = dataOut.StrVal
		case model.Integer:
			dataIn[key] = int(dataOut.IntVal)
		default:
			dataIn[key] = []byte(dataOut.RawValue)
		}
	}

	return &Inject{
		StateImpl: baseState,
		data:      dataIn,
	}, nil
}

type Inject struct {
	StateImpl
	data data.Data[any]
}

func (i *Inject) Run(ctx context.Context, dataIn data.Data[any]) (data.Data[any], error) {
	return dataIn, nil
}
