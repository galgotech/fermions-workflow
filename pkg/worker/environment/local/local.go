package local

import (
	"github.com/serverlessworkflow/sdk-go/v2/model"

	"github.com/galgotech/fermions-workflow/pkg/bus"
	"github.com/galgotech/fermions-workflow/pkg/log"
	"github.com/galgotech/fermions-workflow/pkg/setting"
	"github.com/galgotech/fermions-workflow/pkg/worker/environment"
	"github.com/galgotech/fermions-workflow/pkg/worker/environment/local/event"
	"github.com/galgotech/fermions-workflow/pkg/worker/environment/local/function"
	"github.com/galgotech/fermions-workflow/pkg/worker/environment/local/state"
)

func ProvideNew(setting setting.Setting) environment.NewEnvironment {
	return func() environment.Environment {
		return Provide(setting)
	}
}

func Provide(setting setting.Setting) *Local {
	return &Local{
		log: log.New("worker-environment-local"),
	}
}

type Local struct {
	log       log.Logger
	spec      model.Workflow
	functions environment.MapFunctions
	events    environment.MapEvents
	states    environment.MapStates
}

func (l *Local) Spec() model.Workflow {
	return l.spec
}

func (l *Local) InitializeWorkflow(spec model.Workflow, busEvent bus.Bus) (err error) {
	l.spec = spec
	err = l.initalizeFunctions(spec.Functions)
	if err != nil {
		l.log.Error("environment local fail initialize functions")
		return err
	}

	err = l.initalizeEvents(spec.Events, busEvent)
	if err != nil {
		l.log.Error("environment local fail initialize events")
		return err
	}

	err = l.initializeStates(spec.States)
	if err != nil {
		l.log.Error("environment local fail initialize states")
		return err
	}

	return nil
}

func (l *Local) initalizeFunctions(functionsSpec []model.Function) error {
	functions := make(environment.MapFunctions, 0)
	for _, f := range functionsSpec {
		fPrepared, err := function.New(f)
		if err != nil {
			return err
		}

		functions[f.Name] = fPrepared
	}
	l.functions = functions
	return nil
}

func (l *Local) initalizeEvents(eventsSpec []model.Event, busEvent bus.Bus) error {
	events := make(environment.MapEvents, 0)
	for _, eventSpec := range eventsSpec {
		events[eventSpec.Name] = event.New(eventSpec, busEvent)
	}
	l.events = events
	return nil
}

func (l *Local) initializeStates(statesSpec []model.State) error {
	stateMap := make(environment.MapStates, len(statesSpec))
	for _, stateSpec := range statesSpec {
		t, err := state.New(stateSpec, l.functions, l.events)
		if err != nil {
			l.log.Error("environment local fail new state", "state", stateSpec.Name)
			return err
		}
		stateMap[stateSpec.Name] = t
	}
	l.states = stateMap
	return nil
}

func (l *Local) Start() string {
	return l.spec.Start.StateName
}

func (l *Local) State(name string) environment.State {
	return l.states[name]
}

func (l *Local) CompensateBy(current environment.State) error {
	return nil
}
