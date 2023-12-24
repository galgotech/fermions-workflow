package state

import (
	"context"

	"github.com/serverlessworkflow/sdk-go/v2/model"

	"github.com/galgotech/fermions-workflow/pkg/worker/environment"
	"github.com/galgotech/fermions-workflow/pkg/worker/filter"
)

func newOperation(spec model.OperationState, baseState StateImpl, functions environment.MapFunctions, mapEvents environment.MapEvents) (*Operation, error) {
	actions, err := newAction(spec.Actions, functions, mapEvents)
	if err != nil {
		return nil, err
	}

	return &Operation{
		StateImpl:  baseState,
		actionsLen: len(actions),
		actions:    actions,
	}, nil
}

type Operation struct {
	StateImpl
	actionsLen int
	actions    Actions
	dataFilter filter.Filter
}

func (s *Operation) Run(ctx context.Context, dataIn model.Object) (dataOut model.Object, err error) {
	return s.actions.Run(dataIn)
}
