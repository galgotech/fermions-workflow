package worker

import (
	"sync"
	"time"

	"github.com/serverlessworkflow/sdk-go/v2/model"
	"github.com/serverlessworkflow/sdk-go/v2/parser"

	"github.com/galgotech/fermions-workflow/pkg/bus"
	"github.com/galgotech/fermions-workflow/pkg/log"
	"github.com/galgotech/fermions-workflow/pkg/setting"
	"github.com/galgotech/fermions-workflow/pkg/worker/environment"
	"github.com/galgotech/fermions-workflow/pkg/worker/runtime"
)

func New(s setting.Setting, runtime runtime.Runtime, newEnvironment environment.NewEnvironment, busEvent bus.Bus) (*Worker, error) {
	return &Worker{
		log:            log.New("worker"),
		setting:        s,
		runtime:        runtime,
		newEnvironment: newEnvironment,
		busEvent:       busEvent,
	}, nil
}

type Worker struct {
	log            log.Logger
	setting        setting.Setting
	runtime        runtime.Runtime
	newEnvironment environment.NewEnvironment
	busEvent       bus.Bus
}

func (w *Worker) Execute() error {
	for workflowKey, start := range w.setting.Starts() {
		if start {
			if spec, ok := w.setting.WorkflowSpecs()[workflowKey]; ok {
				err := w.start(spec)
				if err != nil {
					w.log.Error("start workflow", "name", spec.Name, "id", spec.ID, "key", spec.Key, "error", err.Error())
					return err
				}
			}
		}
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for {
			if !w.runtime.WorkflowsRunning() {
				wg.Done()
				break
			}
			time.Sleep(1 * time.Second)
		}
	}()
	wg.Wait()

	return nil
}

func (w *Worker) StartWorkflow(data []byte) error {
	spec, err := parser.FromJSONSource(data)
	if err != nil {
		return err
	}
	return w.start(*spec)
}

func (w *Worker) start(spec model.Workflow) error {
	env := w.newEnvironment()
	err := env.InitializeWorkflow(spec, w.busEvent)
	if err != nil {
		return err
	}
	return w.runtime.Start(env)
}

func (w *Worker) StopWorkflow(wid uint64) {
	w.runtime.Stop(wid)
}

func (w *Worker) Shutdown() {
	w.runtime.Shutdown()
}
