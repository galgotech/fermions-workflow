package state

import (
	"errors"

	"github.com/serverlessworkflow/sdk-go/v2/model"

	"github.com/galgotech/fermions-workflow/pkg/worker/environment"
)

func New(spec model.State, mapFunctions environment.MapFunctions, mapEvents environment.MapEvents) (environment.State, error) {
	stateBase, err := NewBase(spec, mapEvents)
	if err != nil {
		return nil, err
	}

	switch spec.Type {
	case model.StateTypeEvent:
		return newEvent(*spec.EventState, stateBase, mapFunctions, mapEvents)
	case model.StateTypeOperation:
		return newOperation(*spec.OperationState, stateBase, mapFunctions, mapEvents)
	case model.StateTypeSwitch:
		return newSwitch(*spec.SwitchState, stateBase)
	case model.StateTypeInject:
		return newInject(*spec.InjectState, stateBase)
	case model.StateTypeCallback:
		return newCallback(*spec.CallbackState, stateBase, mapFunctions, mapEvents)
	}

	return nil, errors.New("transition not implemented")
}
