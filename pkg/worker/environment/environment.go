package environment

import (
	"context"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/serverlessworkflow/sdk-go/v2/model"

	"github.com/galgotech/fermions-workflow/pkg/bus"
	"github.com/galgotech/fermions-workflow/pkg/worker/data"
)

type MapFunctions map[string]Function
type MapEvents map[string]Event
type MapStates map[string]State

type NewEnvironment func() Environment

type Environment interface {
	Spec() model.Workflow
	InitializeWorkflow(spec model.Workflow, busEvent bus.Bus) error
	Start() (StateStart, error)
	State(name string) State
	CompensateBy(transition State) error
}

type Function interface {
	Init() error
	Run(data data.Data[any]) (data.Data[any], error)
}

type Event interface {
	Produce(ctx context.Context, event cloudevents.Event) error
	Consume(ctx context.Context) (cloudevents.Event, error)
}

type State interface {
	Type() model.StateType
	Name() string

	FilterInput(data.Data[any]) (data.Data[any], error)
	FilterOutput(data.Data[any]) (data.Data[any], error)
	Run(ctx context.Context, dataIn data.Data[any]) (data.Data[any], error)
	ProduceEvents(ctx context.Context, data data.Data[any]) error
	Next(ctx context.Context) <-chan StateStart
}

type StateStart interface {
	Ctx() context.Context
	State() string
}
