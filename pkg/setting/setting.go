package setting

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"

	"github.com/serverlessworkflow/sdk-go/v2/model"
	"github.com/serverlessworkflow/sdk-go/v2/parser"

	"github.com/galgotech/fermions-workflow/pkg/log"
)

func New() *Setting {
	return &Setting{
		log:           log.New("setting"),
		WorkflowSpecs: map[string]model.Workflow{},
		Starts:        map[string]bool{},
	}
}

type Setting struct {
	Bus           Bus
	log           log.Logger
	WorkflowSpecs map[string]model.Workflow
	Starts        map[string]bool
}

func (s *Setting) LoadConfig(confPath string) error {
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

		s.Bus.Redis = url.Value()
	}

	return nil
}

func (s *Setting) ParseWorkflow(filePath string) error {
	info, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	paths := []string{}
	if info.IsDir() {
		files, err := ioutil.ReadDir(filePath)
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

func (s *Setting) AddWorkflow(workflowSpec model.Workflow) {
	workflowKey := ""
	if workflowSpec.Key != "" {
		workflowKey = workflowSpec.Key
	} else {
		workflowKey = workflowSpec.ID
	}

	s.log.Info("add workflowspec", "keyOrId", workflowKey)

	s.Starts[workflowKey] = false
	s.WorkflowSpecs[workflowKey] = workflowSpec
}

func (s *Setting) AddStart(starts []string) {
	for _, start := range starts {
		if _, ok := s.Starts[start]; ok {
			s.Starts[start] = true
		} else {
			s.log.Error("workflowspec not found", "keyOrId", start)
		}
	}
}
