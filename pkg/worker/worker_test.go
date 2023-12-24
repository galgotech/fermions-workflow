package worker

import (
	"testing"

	"github.com/serverlessworkflow/sdk-go/v2/model"
	"github.com/stretchr/testify/assert"

	"github.com/galgotech/fermions-workflow/pkg/bus"
	"github.com/galgotech/fermions-workflow/pkg/setting"
	"github.com/galgotech/fermions-workflow/pkg/test"
	"github.com/galgotech/fermions-workflow/pkg/worker/environment"
	"github.com/galgotech/fermions-workflow/pkg/worker/runtime"
)

func TestExecute(t *testing.T) {
	t.Run("execute", func(t *testing.T) {
		workflowBuilder := model.NewWorkflowBuilder()
		workflow := workflowBuilder.Build()

		worker, fr := setup(workflow)
		worker.Execute()
		assert.Equal(t, 1, fr.Env.(*environmentStub).Count)
		assert.Equal(t, 1, fr.Count)
	})
}

func setup(workflowSpec model.Workflow) (*Worker, *runtimeStub) {
	busEvent := test.NewBusStub()
	fr := &runtimeStub{}
	setting := setting.New().(*setting.FermionsSetting)
	setting.AddWorkflow(workflowSpec)
	setting.AddStart([]string{"test"})
	worker, _ := New(setting, fr, newEnvironment, busEvent)

	return worker, fr
}

func newEnvironment() environment.Environment {
	return &environmentStub{}
}

// =====
type environmentStub struct {
	environment.Environment
	Count int
}

func (s *environmentStub) InitializeWorkflow(spec model.Workflow, busEvent bus.Bus) error {
	s.Count++
	return nil
}

// =====
type runtimeStub struct {
	runtime.Runtime
	Env   environment.Environment
	Count int
}

func (s *runtimeStub) Start(env environment.Environment) error {
	s.Env = env
	s.Count++
	return nil
}
