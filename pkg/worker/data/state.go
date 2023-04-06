package data

import (
	"errors"

	"github.com/galgotech/fermions-workflow/pkg/log"
)

var ErrStateNotFound = errors.New("state not found")

func Provide() Manager {
	return &WorkflowManager{
		log:    log.New("worker-runtime-process"),
		states: make(map[string]Data[any], 0),
	}
}

type Manager interface {
	State(name string) Data[any]
	SetState(name string, state Data[any])
}

type WorkflowManager struct {
	log    log.Logger
	states map[string]Data[any]
}

func (m *WorkflowManager) State(name string) Data[any] {
	if val, ok := m.states[name]; ok {
		return val
	}

	s := Data[any]{}
	m.states[name] = s
	return s
}

func (m *WorkflowManager) SetState(name string, state Data[any]) {
	m.states[name] = state
}
