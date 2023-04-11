package standalone

import (
	"sync"

	"github.com/galgotech/fermions-workflow/pkg/log"
	"github.com/galgotech/fermions-workflow/pkg/server"
	"github.com/galgotech/fermions-workflow/pkg/setting"
	"github.com/galgotech/fermions-workflow/pkg/worker"
)

func New(setting *setting.Setting, server *server.Server, worker *worker.Worker) *Standalone {
	return &Standalone{
		log:     log.New("standalone"),
		setting: setting,
		worker:  worker,
		server:  server,
	}
}

type Standalone struct {
	log     log.Logger
	setting *setting.Setting
	worker  *worker.Worker
	server  *server.Server
}

func (s *Standalone) Execute() (err error) {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		err = s.worker.Execute()
		wg.Done()
	}()

	go func() {
		err = s.server.Execute()
		wg.Done()
	}()

	wg.Wait()
	return
}
