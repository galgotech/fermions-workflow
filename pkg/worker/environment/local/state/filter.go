package state

import (
	"github.com/galgotech/fermions-workflow/pkg/worker/filter"
	"github.com/serverlessworkflow/sdk-go/v2/model"
)

func initializeStateDataFilter(spec *model.StateDataFilter) (filter.Filter, filter.Filter, error) {
	var input, output string
	if spec == nil {
		input = ""
		output = ""
	} else {
		input = spec.Input
		input = spec.Output
	}

	filterInput, err := filter.NewFilter(input)
	if err != nil {
		return nil, nil, err
	}

	filterOutput, err := filter.NewFilter(output)
	if err != nil {
		return nil, nil, err
	}

	return filterInput, filterOutput, nil
}

func initializeEventDataFilter(sepc model.EventDataFilter) (filter.Filter, filter.Filter, error) {
	filterData, err := filter.NewFilter("")
	if err != nil {
		return nil, nil, err
	}

	filterToStateData, err := filter.NewFilter("")
	if err != nil {
		return nil, nil, err
	}
	if sepc.UseData {
		var err error
		filterData, err = filter.NewFilter(sepc.Data)
		if err != nil {
			return nil, nil, err
		}
		filterToStateData, err = filter.NewFilter(sepc.ToStateData)
		if err != nil {
			return nil, nil, err
		}
	}

	return filterData, filterToStateData, nil
}

func initializeActionDataFilter(spec model.ActionDataFilter) (filter.Filter, filter.Filter, filter.Filter, error) {
	filterFromStateData, err := filter.NewFilter(spec.FromStateData)
	if err != nil {
		return nil, nil, nil, err
	}

	var results, toStateData string
	if spec.UseResults {
		results = spec.Results
		toStateData = spec.ToStateData
	} else {
		results = ""
		toStateData = ""
	}

	filterResults, err := filter.NewFilter(results)
	if err != nil {
		return nil, nil, nil, err
	}

	filterToStateData, err := filter.NewFilter(toStateData)
	if err != nil {
		return nil, nil, nil, err
	}

	return filterResults, filterToStateData, filterFromStateData, nil
}
