package standalonewasm

import (
	"github.com/galgotech/fermions-workflow/pkg/log"
	"github.com/galgotech/fermions-workflow/pkg/setting"
	"github.com/galgotech/fermions-workflow/pkg/worker"
)

func New(s setting.Setting, worker *worker.Worker) (*StandaloneWasm, error) {
	return &StandaloneWasm{
		log:     log.New("standaloneWasm"),
		setting: s,
		worker:  worker,
	}, nil
}

type StandaloneWasm struct {
	log     log.Logger
	setting setting.Setting
	worker  *worker.Worker
}

func (s *StandaloneWasm) Execute(spec []byte) error {
	return s.worker.StartWorkflow(spec)
}

func (s *StandaloneWasm) Kill(id string) error {
	return s.worker.StopWorkflow(id)
}
