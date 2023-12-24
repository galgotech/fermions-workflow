package environment

import (
	"context"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/serverlessworkflow/sdk-go/v2/model"

	"github.com/galgotech/fermions-workflow/pkg/bus"
)

type MapFunctions map[string]Function
type MapEvents map[string]Event
type MapStates map[string]State

type NewEnvironment func() Environment

type Environment interface {
	Spec() model.Workflow
	InitializeWorkflow(spec model.Workflow, busEvent bus.Bus) error
	Start() string
	State(name string) State
	CompensateBy(transition State) error
}

type Function interface {
	Init() error
	Run(data model.Object) (model.Object, error)
}

type Event interface {
	Produce(ctx context.Context, data model.Object) error
	Consume(ctx context.Context) (cloudevents.Event, error)
}

type State interface {
	Type() model.StateType
	Name() string

	FilterInput(model.Object) (model.Object, error)
	FilterOutput(model.Object) (model.Object, error)
	Run(ctx context.Context, dataIn model.Object) (model.Object, error)
	ProduceEvents(ctx context.Context, data model.Object) error
	Next() (string, bool)
}
