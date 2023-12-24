package setting

import (
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"

	"github.com/serverlessworkflow/sdk-go/v2/model"
	"github.com/serverlessworkflow/sdk-go/v2/parser"

	"github.com/galgotech/fermions-workflow/pkg/log"
)

func New() Setting {
	return &FermionsSetting{
		log:           log.New("setting"),
		workflowSpecs: map[string]model.Workflow{},
		starts:        map[string]bool{},
	}
}

type Setting interface {
	Bus() Bus
	WorkflowSpecs() map[string]model.Workflow
	Starts() map[string]bool
}

type FermionsSetting struct {
	bus           Bus
	log           log.Logger
	workflowSpecs map[string]model.Workflow
	starts        map[string]bool
}

func (s *FermionsSetting) Bus() Bus {
	return s.bus
}

func (s *FermionsSetting) WorkflowSpecs() map[string]model.Workflow {
	return s.workflowSpecs
}

func (s *FermionsSetting) Starts() map[string]bool {
	return s.starts
}

func (s *FermionsSetting) LoadConfig(confPath string) error {
	s.log.Info("load config", "path", confPath)

	cfg, err := ini.Load(confPath)
	if err != nil {
		return err
	}

	if cfg.HasSection("redis") {
		section := cfg.Section("redis")
		url, err := section.GetKey("url")
		if err != nil {
			return err
		}

		s.bus.Redis = url.Value()
	}

	return nil
}

func (s *FermionsSetting) ParseWorkflow(filePath string) error {
	info, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	paths := []string{}
	if info.IsDir() {
		files, err := os.ReadDir(filePath)
		if err != nil {
			return err
		}
		for _, f := range files {
			paths = append(paths, filepath.Join(filePath, f.Name()))
		}
	} else {
		paths = append(paths, filePath)
	}

	for _, path := range paths {
		workflowSpec, err := parser.FromFile(path)
		if err != nil {
			return err
		}
		s.AddWorkflow(*workflowSpec)
	}

	return nil
}

func (s *FermionsSetting) AddWorkflow(workflowSpec model.Workflow) {
	workflowKey := ""
	if workflowSpec.Key != "" {
		workflowKey = workflowSpec.Key
	} else {
		workflowKey = workflowSpec.ID
	}

	s.log.Info("add workflowspec", "keyOrId", workflowKey)

	s.starts[workflowKey] = false
	s.workflowSpecs[workflowKey] = workflowSpec
}

func (s *FermionsSetting) AddStart(starts []string) {
	for _, start := range starts {
		if _, ok := s.starts[start]; ok {
			s.starts[start] = true
		} else {
			s.log.Error("workflowspec not found", "keyOrId", start)
		}
	}
}
