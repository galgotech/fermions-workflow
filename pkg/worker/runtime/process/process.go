package process

import (
	"context"

	"github.com/google/uuid"
	"github.com/serverlessworkflow/sdk-go/v2/model"

	"github.com/galgotech/fermions-workflow/pkg/bus"
	"github.com/galgotech/fermions-workflow/pkg/log"
	"github.com/galgotech/fermions-workflow/pkg/worker/data"
	"github.com/galgotech/fermions-workflow/pkg/worker/environment"
)

type ContextKey string

const (
	ContextKeyTrace ContextKey = "trace"
)

func Provide(dataManager data.Manager, workflowBus bus.Bus) *WorkflowRuntime {
	runtime := &WorkflowRuntime{
		log: log.New("worker-runtime-process"),

		dataManager:      dataManager,
		bus:              workflowBus,
		workflowsRunning: make(map[uint64]*stateIn),
	}

	return runtime
}

type WorkflowRuntime struct {
	log log.Logger

	dataManager      data.Manager
	bus              bus.Bus
	currentPid       uint64
	workflowsRunning map[uint64]*stateIn
}

func (r *WorkflowRuntime) Start(env environment.Environment) error {
	stateCtx, err := newStateIn(env)
	if err != nil {
		return err
	}

	r.workflowsRunning[r.currentPid] = stateCtx
	stateCtx.wid = r.currentPid
	r.currentPid += 1

	r.state(stateCtx)
	return nil
}

func (r *WorkflowRuntime) Shutdown() {
	r.log.Info("shutdown")
	for wid := range r.workflowsRunning {
		r.Stop(wid)
	}
}

func (r *WorkflowRuntime) Stop(wid uint64) {
	r.log.Info("stop workflow", "wid", wid)
	if execState, ok := r.workflowsRunning[wid]; ok {
		// TODO: Check cancel is corret closing go routines
		execState.cancel()
		delete(r.workflowsRunning, wid)
	}
}

func (r *WorkflowRuntime) WorkflowsRunning() bool {
	// TODO: Check concurrency from workflowsRunning
	return len(r.workflowsRunning) > 0
}

func (r *WorkflowRuntime) state(stateCtx *stateIn) {
	go func() {
		ctx := stateCtx.ctx
		env := stateCtx.env
		state := stateCtx.env.State(stateCtx.state)

		r.log.Info("state run", "workflow", env.Spec().Name, "state", stateCtx.state, "type", state.Type(), "trace", ctx.Value(ContextKeyTrace))

		dataIn := r.dataManager.Get()

		r.log.Debug("state filter input", "workflow", env.Spec().Name, "state", stateCtx.state, "trace", ctx.Value(ContextKeyTrace))
		dataIn, err := state.FilterInput(dataIn)
		if err != nil {
			r.log.Error("state error", "state", state.Name(), "err", err.Error(), "trace", ctx.Value(ContextKeyTrace))
			return
		}

		r.log.Debug("state run", "workflow", env.Spec().Name, "state", stateCtx.state, "trace", ctx.Value(ContextKeyTrace))
		dataOut, err := state.Run(ctx, dataIn)
		if err != nil {
			r.log.Error("state error", "state", state.Name(), "err", err.Error(), "trace", ctx.Value(ContextKeyTrace))
			return
		}

		r.log.Debug("state output", "workflow", env.Spec().Name, "state", stateCtx.state, "trace", ctx.Value(ContextKeyTrace))
		dataOut, err = state.FilterOutput(dataOut)
		if err != nil {
			r.log.Error("state error", "state", state.Name(), "err", err.Error(), "trace", ctx.Value(ContextKeyTrace))
			return
		}

		r.log.Debug("state data save", "workflow", env.Spec().Name, "state", stateCtx.state, "trace", ctx.Value(ContextKeyTrace))
		err = r.dataManager.Set(dataOut)
		if err != nil {
			r.log.Error("state set error", "state", state.Name(), "err", err.Error(), "trace", ctx.Value(ContextKeyTrace))
			return
		}

		r.log.Debug("state compensate by", "workflow", env.Spec().Name, "state", stateCtx.state, "trace", ctx.Value(ContextKeyTrace))
		err = env.CompensateBy(state)
		if err != nil {
			r.log.Error("state error", "state", state.Name(), "err", err.Error(), "trace", ctx.Value(ContextKeyTrace))
			return
		}

		r.log.Debug("state produce events", "workflow", env.Spec().Name, "state", stateCtx.state, "trace", ctx.Value(ContextKeyTrace))
		err = state.ProduceEvents(ctx, dataOut)
		if err != nil {
			r.log.Error("produce events", "state", state.Name(), "err", err.Error(), "trace", ctx.Value(ContextKeyTrace))
			return
		}

		// Start a new EventState to waiting a new workflow execution
		if state.Type() == model.StateTypeEvent && env.Start() == stateCtx.state {
			r.Start(env)
		}

		// Next state
		stateName, ok := state.Next()
		if !ok {
			r.log.Debug("workflow done", "currentState", state.Name(), "trace", ctx.Value(ContextKeyTrace))
			r.Stop(stateCtx.wid)
			return
		}
		stateCtx.state = stateName
		r.state(stateCtx)
	}()
}

func newStateIn(env environment.Environment) (*stateIn, error) {
	trace, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	ctx = context.WithValue(ctx, ContextKeyTrace, trace.String())
	return &stateIn{
		ctx:    ctx,
		cancel: cancel,
		env:    env,
		state:  env.Start(),
	}, nil
}

type stateIn struct {
	wid    uint64
	ctx    context.Context
	cancel context.CancelFunc

	env   environment.Environment
	state string
}
