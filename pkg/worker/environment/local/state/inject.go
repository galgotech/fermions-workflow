package state

import (
	"context"

	"github.com/serverlessworkflow/sdk-go/v2/model"
)

func newInject(spec model.InjectState, baseState StateImpl) (*Inject, error) {
	return &Inject{
		StateImpl: baseState,
		data:      spec.Data,
	}, nil
}

type Inject struct {
	StateImpl
	data map[string]model.Object
}

func (i *Inject) Run(ctx context.Context, dataIn model.Object) (model.Object, error) {
	return dataIn, nil
}
