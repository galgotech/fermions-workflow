package runtime

import (
	"context"

	"github.com/galgotech/fermions-workflow/pkg/worker/environment"
)

type Runtime interface {
	Start(ctx context.Context, env environment.Environment) error
	Shutdown()
}
