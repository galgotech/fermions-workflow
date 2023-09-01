package state

import (
	"github.com/serverlessworkflow/sdk-go/v2/model"

	"github.com/galgotech/fermions-workflow/pkg/worker/data"
	"github.com/galgotech/fermions-workflow/pkg/worker/environment"
	"github.com/galgotech/fermions-workflow/pkg/worker/filter"
)

func newAction(specs []model.Action, functions environment.MapFunctions) (Actions, error) {
	actions := make(Actions, len(specs))
	for i, spec := range specs {
		filterFromStateData, filterResults, filterToStateData, err := initializeActionDataFilter(spec.ActionDataFilter)
		if err != nil {
			return nil, err
		}

		actions[i] = &Action{
			Function:            functions[spec.FunctionRef.RefName],
			filterFromStateData: filterFromStateData,
			filterResults:       filterResults,
			filterToStateData:   filterToStateData,
		}
	}

	return actions, nil
}

type Action struct {
	Function environment.Function
	// Workflow expression that filters state data that can be used by the action
	filterFromStateData filter.Filter
	// Workflow expression that filters the actions data results
	filterResults filter.Filter
	// Workflow expression that selects a state data element to which the action results should be added/merged into. If not specified denotes the top-level state data element
	filterToStateData filter.Filter
}

func (a *Action) Run(dataIn model.Object) (dataOut model.Object, err error) {
	dataIn, err = a.filterFromStateData.Run(dataIn)
	if err != nil {
		return
	}

	dataOut, err = a.Function.Run(dataIn)
	if err != nil {
		return
	}

	dataOut, err = a.filterResults.Run(dataOut)
	if err != nil {
		return
	}

	dataOut, err = a.filterToStateData.Run(dataOut)
	if err != nil {
		return
	}

	return
}

type Actions []*Action

func (a Actions) Run(dataIn model.Object) (dataOut model.Object, err error) {
	for _, action := range a {
		dataOut, err = action.Run(dataIn)
		if err != nil {
			return data.ObjectNil, err
		}
		dataIn = dataOut
	}

	return
}
