package state

import (
	"context"

	"github.com/serverlessworkflow/sdk-go/v2/model"

	"github.com/galgotech/fermions-workflow/pkg/worker/environment"
)

func newSwitch(spec model.SwitchState, stateBase StateImpl) (environment.State, error) {

	// conditions []SwitchCondition, defaultCondition SwitchCondition

	conditions := make([]SwitchCondition, len(spec.DataConditions))
	for i, dataCondition := range spec.DataConditions {
		t := SwitchCondition{condition: dataCondition.Condition, toState: dataCondition.Transition.NextState}
		conditions[i] = t
	}

	defaultCondition := SwitchCondition{condition: "", toState: spec.DefaultCondition.Transition.NextState}

	return &Switch{
		StateImpl:        stateBase,
		conditions:       conditions,
		defaultCondition: defaultCondition,
	}, nil
}

type Switch struct {
	StateImpl
	conditions       []SwitchCondition
	defaultCondition SwitchCondition
}

func (s *Switch) Run(ctx context.Context, dataIn model.Object) (model.Object, error) {
	return dataIn, nil
}

type SwitchCondition struct {
	condition string
	toState   string
}
