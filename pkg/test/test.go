package test

import "github.com/serverlessworkflow/sdk-go/v2/model"

var Events = []model.Event{
	{
		Name: "testEvent",
	},
}

var Functions = []model.Function{
	{
		Name:      "test",
		Type:      "rest",
		Operation: "https://galgo.tech/functions.json#test",
	},
}

var Action1 = model.Action{
	Name: "action1",
	FunctionRef: &model.FunctionRef{
		RefName: "test",
	},
	ActionDataFilter: model.ActionDataFilter{
		FromStateData: "",
		UseResults:    true,
		Results:       "",
		ToStateData:   "",
	},
}

var OnEvent = model.OnEvents{
	EventRefs: []string{"nameEvent0"},
	Actions:   []model.Action{Action1},
	// EventDataFilter:
}

var EventState = model.State{
	BaseState: model.BaseState{
		Type: model.StateTypeEvent,
		Name: "event1",
		End: &model.End{
			Terminate: true,
		},
	},
	EventState: &model.EventState{
		OnEvents: []model.OnEvents{OnEvent},
	},
}

var StateInject = model.State{
	BaseState: model.BaseState{
		Type: model.StateTypeOperation,
		Name: "stateInject",
		End: &model.End{
			Terminate: true,
		},
	},
	InjectState: &model.InjectState{
		Data: map[string]model.Object{
			"test0": model.FromString("testValStr"),
			"test1": model.FromInt(1),
			"test2": model.FromRaw("bytes"),
		},
	},
}

var CallbackState = model.State{
	BaseState: model.BaseState{
		Type: model.StateTypeOperation,
		Name: "stateCallback",
		End: &model.End{
			Terminate: true,
		},
	},
	CallbackState: &model.CallbackState{
		Action: Action1,
	},
}

var States = []model.State{
	{
		BaseState: model.BaseState{
			Type: model.StateTypeOperation,
			Name: "stateStart",
			Transition: &model.Transition{
				NextState: "state2",
			},
		},
		OperationState: &model.OperationState{
			ActionMode: model.ActionModeSequential,
			Actions:    []model.Action{Action1},
		},
	},
	{
		BaseState: model.BaseState{
			Type: model.StateTypeOperation,
			Name: "state2",
			End: &model.End{
				Terminate: true,
			},
		},
		OperationState: &model.OperationState{
			ActionMode: model.ActionModeSequential,
			Actions: []model.Action{{
				Name: "action1",
			}},
		},
	},
}

var Workflow = model.Workflow{
	BaseWorkflow: model.BaseWorkflow{
		ID:          "test",
		Name:        "test",
		Description: "test",
		Version:     "0.0.0",
		SpecVersion: "0.8",
		Start: &model.Start{
			StateName: "stateStart",
		},
	},
	Events:    Events,
	Functions: Functions,
	States:    States,
}
