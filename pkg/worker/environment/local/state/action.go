package state

import (
	"time"

	"github.com/serverlessworkflow/sdk-go/v2/model"

	"github.com/galgotech/fermions-workflow/pkg/worker/data"
	"github.com/galgotech/fermions-workflow/pkg/worker/environment"
	"github.com/galgotech/fermions-workflow/pkg/worker/filter"
)

func newAction(specs []model.Action, functions environment.MapFunctions, events environment.MapEvents) (Actions, error) {
	actions := make(Actions, len(specs))
	for i, spec := range specs {
		filterFromStateData, filterResults, filterToStateData, err := initializeActionDataFilter(spec.ActionDataFilter)
		if err != nil {
			return nil, err
		}

		var sleepBefore time.Duration
		var sleepAfter time.Duration
		if spec.Sleep != nil {
			sleepBefore, err = calcTimeDuration(spec.Sleep.Before)
			if err != nil {
				return nil, err
			}

			sleepAfter, err = calcTimeDuration(spec.Sleep.After)
			if err != nil {
				return nil, err
			}
		}

		actions[i] = &Action{
			spec:     spec,
			Function: functions[spec.FunctionRef.RefName],
			// Event:               events[spec.EventRef],
			sleepBefore:         sleepBefore,
			sleepAfter:          sleepAfter,
			filterFromStateData: filterFromStateData,
			filterResults:       filterResults,
			filterToStateData:   filterToStateData,
		}
	}

	return actions, nil
}

type Action struct {
	spec                model.Action
	Function            environment.Function
	Event               environment.Event
	sleepBefore         time.Duration
	sleepAfter          time.Duration
	filterFromStateData filter.Filter
	filterResults       filter.Filter
	filterToStateData   filter.Filter
}

func (a *Action) Run(dataIn model.Object) (model.Object, error) {
	// Workflow expression that filters state data that can be used by the action
	dataIn, err := a.filterFromStateData.Run(dataIn)
	if err != nil {
		return data.ObjectNil, err
	}

	// Defines time periods workflow execution should sleep before function execution
	if a.sleepBefore.Nanoseconds() > 0 {
		time.Sleep(a.sleepBefore)
	}

	dataOut, err := a.Function.Run(dataIn)
	if err != nil {
		return data.ObjectNil, err
	}

	// Defines time periods workflow execution should sleep after function execution
	if a.sleepAfter.Nanoseconds() > 0 {
		time.Sleep(a.sleepAfter)
	}

	// If set to false, action data results are not added/merged to state data.
	if a.spec.ActionDataFilter.UseResults {
		// Workflow expression that filters the actions data results
		dataOut, err = a.filterResults.Run(dataOut)
		if err != nil {
			return data.ObjectNil, err
		}

		// Workflow expression that selects a state data element to which the action results should be added/merged into. If not specified denotes the top-level state data element
		dataOut, err = a.filterToStateData.Run(dataOut)
		if err != nil {
			return data.ObjectNil, err
		}
		return dataOut, nil
	}

	return data.ObjectNil, nil
}

type Actions []*Action

func (a Actions) Run(dataIn model.Object) (model.Object, error) {
	// dataMerge := model.FromMap(map[string]any{})
	for _, action := range a {
		dataOut, err := action.Run(dataIn)
		if err != nil {
			return data.ObjectNil, err
		}

		dataIn, err = data.Merge(dataIn, dataOut)
		if err != nil {
			return data.ObjectNil, err
		}
	}

	return dataIn, nil
}
