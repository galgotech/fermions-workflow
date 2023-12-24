package data

import (
	"errors"

	"github.com/galgotech/fermions-workflow/pkg/log"
	"github.com/serverlessworkflow/sdk-go/v2/model"
)

var ErrStateNotFound = errors.New("state not found")

func Provide() Manager {
	return &WorkflowManager{
		log:   log.New("worker-runtime-process"),
		state: model.FromMap(map[string]any{}),
	}
}

type Manager interface {
	Get() model.Object
	Set(state model.Object) error
}

type WorkflowManager struct {
	log   log.Logger
	state model.Object
}

func (m *WorkflowManager) Get() model.Object {
	return m.state
}

func (m *WorkflowManager) Set(state model.Object) error {
	if state.Type != model.Map {
		return errors.New("state set require type model.Map")
	}

	if len(m.state.MapValue) == 0 {
		m.state = state
		return nil
	}

	state, err := Merge(m.state, state)
	if err != nil {
		return err
	}
	m.state = state
	return nil
}
