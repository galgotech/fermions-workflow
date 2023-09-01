package data

import (
	"errors"

	"github.com/galgotech/fermions-workflow/pkg/log"
	"github.com/serverlessworkflow/sdk-go/v2/model"
)

var ErrStateNotFound = errors.New("state not found")

func Provide() Manager {
	return &WorkflowManager{
		log:    log.New("worker-runtime-process"),
		states: make(map[string]model.Object, 0),
	}
}

type Manager interface {
	State(name string) model.Object
	SetState(name string, state model.Object)
}

type WorkflowManager struct {
	log    log.Logger
	states map[string]model.Object
}

func (m *WorkflowManager) State(name string) model.Object {
	if val, ok := m.states[name]; ok {
		return val
	}

	s := model.Object{}
	m.states[name] = s
	return s
}

func (m *WorkflowManager) SetState(name string, state model.Object) {
	m.states[name] = state
}
