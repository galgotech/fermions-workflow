//go:build wireinject
// +build wireinject

package standalonewasm

import (
	"github.com/galgotech/fermions-workflow/pkg/bus"
	"github.com/galgotech/fermions-workflow/pkg/setting"
	"github.com/galgotech/fermions-workflow/pkg/worker"
	"github.com/galgotech/fermions-workflow/pkg/worker/data"
	"github.com/galgotech/fermions-workflow/pkg/worker/environment"
	"github.com/galgotech/fermions-workflow/pkg/worker/environment/local"
	"github.com/galgotech/fermions-workflow/pkg/worker/runtime"
	"github.com/galgotech/fermions-workflow/pkg/worker/runtime/process"
	"github.com/google/wire"
)

var wireBasicSet = wire.NewSet(
	local.ProvideNew,
	local.Provide,
	wire.Bind(new(environment.Environment), new(*local.Local)),
	data.Provide,
	process.Provide,
	bus.Provide,
	worker.New,
	New,
	wire.Bind(new(runtime.Runtime), new(*process.WorkflowRuntime)),
)

func Initialize(setting.Setting) (*StandaloneWasm, error) {
	wire.Build(wireBasicSet)
	return &StandaloneWasm{}, nil
}
