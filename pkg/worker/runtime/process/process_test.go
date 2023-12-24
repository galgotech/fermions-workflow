package process

import (
	"context"
	"testing"
	"time"

	"github.com/serverlessworkflow/sdk-go/v2/model"
	"github.com/stretchr/testify/assert"

	"github.com/galgotech/fermions-workflow/pkg/bus"
	"github.com/galgotech/fermions-workflow/pkg/test"
	"github.com/galgotech/fermions-workflow/pkg/worker/data"
	"github.com/galgotech/fermions-workflow/pkg/worker/environment"
)

func TestStart(t *testing.T) {
	dataManager := data.Provide()
	process := Provide(dataManager, test.NewBusStub())

	t.Run("start state operation", func(t *testing.T) {
		env := &environmentStub{}
		process.Start(env)
		time.Sleep(1 * time.Second) // TODO: remove the sleep to improve this test

		assert.Equal(t, 2, env.CountCompensateBy)
		assert.Equal(t, 2, env.CountProduceEvents)
		dataState := dataManager.Get()
		assert.Equal(t, model.FromMap(map[string]any{"test": "test0"}), dataState)
		dataState = dataManager.Get()
		assert.Equal(t, model.FromMap(map[string]any{"test": "test1"}), dataState)
	})
}

// Environment
type environmentStub struct {
	environment.Environment
	CountCompensateBy  int
	CountProduceEvents int
}

func (e *environmentStub) Start() string {
	return "test0"
}

func (s *environmentStub) State(name string) environment.State {
	return &stateEventStub{name: name}
}

func (e *environmentStub) CompensateBy(transition environment.State) error {
	e.CountCompensateBy++
	return nil
}

func (e *environmentStub) ProduceEvents(ctx context.Context, busEvent bus.Bus, transition environment.State, dataIn model.Object) {
	e.CountProduceEvents++
}

// StateOperation
type stateOperationStub struct {
	environment.State
	name string
}

func (s *stateOperationStub) Type() model.StateType {
	return model.StateTypeOperation
}

func (s *stateOperationStub) Name() string {
	return s.name
}

func (s *stateOperationStub) FilterInput(dataIn model.Object) (model.Object, error) {
	return dataIn, nil
}

func (s *stateOperationStub) FilterOutput(dataIn model.Object) (model.Object, error) {
	return dataIn, nil
}

func (s *stateOperationStub) Run(ctx context.Context, dataIn model.Object) (model.Object, error) {
	dataOut := model.FromInterface(map[string]any{"test": s.name})
	return dataOut, nil
}

func (s *stateOperationStub) Next() (string, bool) {
	// ch := make(chan environment.StateStart)
	// go func() {
	// 	defer close(ch)
	// 	// ch <- "test1"
	// }()
	return "test1", true
}

// StateEvent
type stateEventStub struct {
	environment.State
	name string
}

func (s *stateEventStub) Type() model.StateType {
	return model.StateTypeEvent
}

func (s *stateEventStub) Name() string {
	return s.name
}

func (s *stateEventStub) FilterInput(dataIn model.Object) (model.Object, error) {
	return dataIn, nil
}

func (s *stateEventStub) FilterOutput(dataIn model.Object) (model.Object, error) {
	return dataIn, nil
}

func (s *stateEventStub) Run(ctx context.Context, dataIn model.Object) (model.Object, error) {
	dataOut := model.FromInterface(map[string]any{"test": s.name})
	return dataOut, nil
}

func (s *stateEventStub) Next() (string, bool) {
	// ch := make(chan environment.StateStart)
	// go func() {
	// 	defer close(ch)
	// }()
	return "", false
}
