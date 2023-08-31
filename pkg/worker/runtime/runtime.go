package runtime

import (
	"github.com/galgotech/fermions-workflow/pkg/worker/environment"
)

type Runtime interface {
	Start(env environment.Environment) error
	WorkflowsRunning() bool
	Stop(id uint64)
	Shutdown()
}
