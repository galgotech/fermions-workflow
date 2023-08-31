package cmd

import (
	"os"

	"github.com/urfave/cli/v2"

	"github.com/galgotech/fermions-workflow/pkg/setting"
	"github.com/galgotech/fermions-workflow/pkg/worker"
)

func Worker() error {
	setting := setting.New().(*setting.FermionsSetting)

	app := &cli.App{
		Name:    "fermions-workflow-worker",
		Usage:   "Workflow worker",
		Authors: authors,
		Flags:   globalFlags(setting, true, true),
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "",
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name:     "exec",
						Usage:    "exec workflow",
						Required: true,
						Action: func(c *cli.Context, exec []string) error {
							setting.AddStart(exec)
							return nil
						},
					},
				},
				Action: func(c *cli.Context) error {
					fermionsWorker, err := worker.Initialize(setting)
					if err != nil {
						return err
					}
					return fermionsWorker.Execute()
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		return err
	}

	return nil
}
