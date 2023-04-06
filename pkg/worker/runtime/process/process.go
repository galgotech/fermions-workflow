package process

import (
	"context"
	"time"

	"github.com/galgotech/fermions-workflow/pkg/bus"
	"github.com/galgotech/fermions-workflow/pkg/concurrency"
	"github.com/galgotech/fermions-workflow/pkg/log"
	"github.com/galgotech/fermions-workflow/pkg/worker/data"
	"github.com/galgotech/fermions-workflow/pkg/worker/environment"
)

func Provide(dataManager data.Manager, workflowBus bus.Bus) *WorkflowRuntime {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	runtime := &WorkflowRuntime{
		log:         log.New("worker-runtime-process"),
		ctx:         ctx,
		cancel:      cancel,
		dataManager: dataManager,
		bus:         workflowBus,
		stateIn:     make(chan stateIn),
	}

	runtime.init()

	return runtime
}

func newStateIn(ctx context.Context, env environment.Environment, state environment.State) stateIn {
	return stateIn{
		ctx:   ctx,
		env:   env,
		state: state,
	}
}

type stateIn struct {
	ctx   context.Context
	env   environment.Environment
	state environment.State
}

type WorkflowRuntime struct {
	log log.Logger

	ctx         context.Context
	cancel      context.CancelFunc
	dataManager data.Manager
	bus         bus.Bus

	stateIn chan stateIn
}

func (r *WorkflowRuntime) init() {
	go func() {
		for state := range concurrency.OrDoneCtx(r.ctx, r.stateIn) {
			go func(state stateIn) {
				r.state(state.ctx, state.env, state.state)
			}(state)
		}
	}()
}

func (r *WorkflowRuntime) Shutdown() {
	r.log.Info("shutdown")
	close(r.stateIn)
	r.cancel()
}

func (r *WorkflowRuntime) Start(ctx context.Context, env environment.Environment) error {

	stateGenerator, err := env.Start()
	if err != nil {
		return err
	}

	stateIn := newStateIn(stateGenerator.Ctx(), env, env.State(stateGenerator.State()))
	go func() {
		select {
		case <-time.After(5 * time.Second):
			r.log.Error("start workflow timeout", "name", env.Spec().Name)
		case r.stateIn <- stateIn:
		}
	}()

	return nil
}

func (r *WorkflowRuntime) state(ctx context.Context, env environment.Environment, state environment.State) {
	go func() {
		r.log.Info("state start", "name", env.Spec().Name, "type", state.Type(), "trace", ctx.Value("trace"))

		var dataIn data.Data[any]
		var dataOut data.Data[any]

		r.log.Debug("state filter input", "name", env.Spec().Name, "trace", ctx.Value("trace"))
		dataIn, err := state.FilterInput(dataIn)
		if err != nil {
			r.log.Error("state error", "state", state.Name(), "err", err.Error(), "trace", ctx.Value("trace"))
			return
		}

		r.log.Debug("state run", "name", env.Spec().Name, "trace", ctx.Value("trace"))
		dataOut, err = state.Run(ctx, dataIn)
		if err != nil {
			r.log.Error("state error", "state", state.Name(), "err", err.Error(), "trace", ctx.Value("trace"))
			return
		}

		r.log.Debug("state output", "name", env.Spec().Name, "trace", ctx.Value("trace"))
		dataOut, err = state.FilterOutput(dataOut)
		if err != nil {
			r.log.Error("state error", "state", state.Name(), "err", err.Error(), "trace", ctx.Value("trace"))
			return
		}

		r.log.Debug("state data save", "name", env.Spec().Name, "trace", ctx.Value("trace"))
		r.dataManager.SetState(state.Name(), dataOut)

		r.log.Debug("state compensate by", "name", env.Spec().Name, "trace", ctx.Value("trace"))
		err = env.CompensateBy(state)
		if err != nil {
			r.log.Error("state error", "state", state.Name(), "err", err.Error(), "trace", ctx.Value("trace"))
			return
		}

		r.log.Debug("state produce events", "name", env.Spec().Name, "trace", ctx.Value("trace"))
		err = state.ProduceEvents(ctx, dataOut)
		if err != nil {
			r.log.Error("produce events", "state", state.Name(), "err", err.Error(), "trace", ctx.Value("trace"))
			return
		}

		for nextState := range concurrency.OrDoneCtx(ctx, state.Next(ctx)) {
			r.log.Debug("state next", "currentState", state.Name(), "nextState", nextState.State(), "trace", ctx.Value("trace"))
			r.stateIn <- newStateIn(nextState.Ctx(), env, env.State(nextState.State()))
		}
	}()
}
