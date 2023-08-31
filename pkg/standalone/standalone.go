package standalone

import (
	"context"
	"os/signal"
	"sync"
	"syscall"

	"github.com/galgotech/fermions-workflow/pkg/log"
	"github.com/galgotech/fermions-workflow/pkg/server"
	"github.com/galgotech/fermions-workflow/pkg/setting"
	"github.com/galgotech/fermions-workflow/pkg/worker"
)

func New(s setting.Setting, server *server.Server, worker *worker.Worker) *Standalone {
	return &Standalone{
		log:     log.New("standalone"),
		setting: s,
		worker:  worker,
		server:  server,
	}
}

type Standalone struct {
	log     log.Logger
	setting setting.Setting
	worker  *worker.Worker
	server  *server.Server
}

func (s *Standalone) Execute() {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		// run server
		go func() {
			if err := s.worker.Execute(); err != nil {
				s.log.Error("server", "error", err)
				return
			}
		}()

		// run worker
		go func() {
			if err := s.server.Execute(); err != nil {
				s.log.Error("server", "error", err)
				return
			}
		}()

		// Create context that listens for the interrupt signal from the OS.
		ctx := context.Background()
		ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)

		// Listen for the interrupt signal.
		<-ctx.Done()

		// TODO: implement force shutting down "press Ctrl+C again to force"
		s.log.Info("shutting down gracefully")
		stop()

		s.worker.Shutdown()
		s.log.Info("Exiting")

		wg.Done()
	}()

	wg.Wait()
}
