package standalonewasm

import (
	"github.com/serverlessworkflow/sdk-go/v2/model"

	"github.com/galgotech/fermions-workflow/pkg/log"
	"github.com/galgotech/fermions-workflow/pkg/setting"
)

func NewSetting() setting.Setting {
	return &SettingWasm{
		log: log.New("setting"),
	}
}

type SettingWasm struct {
	bus setting.Bus
	log log.Logger
}

func (s *SettingWasm) Bus() setting.Bus {
	return s.bus
}

func (s *SettingWasm) WorkflowSpecs() map[string]model.Workflow {
	return map[string]model.Workflow{}
}

func (s *SettingWasm) Starts() map[string]bool {
	return map[string]bool{}
}
