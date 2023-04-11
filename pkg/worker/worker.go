package worker

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/serverlessworkflow/sdk-go/v2/model"

	"github.com/galgotech/fermions-workflow/pkg/bus"
	"github.com/galgotech/fermions-workflow/pkg/log"
	"github.com/galgotech/fermions-workflow/pkg/setting"
	"github.com/galgotech/fermions-workflow/pkg/worker/environment"
	"github.com/galgotech/fermions-workflow/pkg/worker/runtime"
)

func New(setting *setting.Setting, runtime runtime.Runtime, newEnvironment environment.NewEnvironment, busEvent bus.Bus) *Worker {
	return &Worker{
		log:            log.New("worker"),
		setting:        setting,
		runtime:        runtime,
		newEnvironment: newEnvironment,
		busEvent:       busEvent,
	}
}

type Worker struct {
	log            log.Logger
	setting        *setting.Setting
	runtime        runtime.Runtime
	newEnvironment environment.NewEnvironment
	busEvent       bus.Bus
}

func (w *Worker) Execute() error {
	for workFlowKey, start := range w.setting.Starts {
		if start {
			if spec, ok := w.setting.WorkflowSpecs[workFlowKey]; ok {
				err := w.start(spec)
				if err != nil {
					w.log.Error("workflowSpec start", "name", spec.Name, "id", spec.ID, "key", spec.Key, "error", err.Error())
				}
			}
		}
	}

	// Create context that listens for the interrupt signal from the OS.
	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)

	// Listen for the interrupt signal.
	<-ctx.Done()
	fmt.Println("")

	// TODO: implement force shutting down "press Ctrl+C again to force"
	w.log.Info("shutting down gracefully")
	stop()

	w.runtime.Shutdown()
	w.log.Info("Exiting")
	return nil
}

func (w *Worker) start(spec model.Workflow) error {
	env := w.newEnvironment()
	if err := env.InitializeWorkflow(spec, w.busEvent); err != nil {
		return err
	}

	ctx := context.Background()
	return w.runtime.Start(ctx, env)
}
