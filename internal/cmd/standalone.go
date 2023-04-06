package cmd

import (
	"os"

	"github.com/urfave/cli/v2"

	"github.com/galgotech/fermions-workflow/pkg/setting"
	"github.com/galgotech/fermions-workflow/pkg/worker"
)

func WorkerStandalone() error {

	setting := setting.New()
	wokerStandalone, err := worker.Initialize(setting)
	if err != nil {
		return err
	}

	app := &cli.App{
		Name:  "workflow-standalone",
		Usage: "Workflow runtime standalone",
		Authors: []*cli.Author{
			{
				Name:  "GalgoTech",
				Email: "andre@galgo.tech",
			},
		},
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:        "conf",
				Usage:       "",
				DefaultText: "conf/default.ini",
				Required:    false,
				Action: func(c *cli.Context, path cli.Path) error {
					setting.ConfPath = path
					return nil
				},
			},
			&cli.PathFlag{
				Name:     "workflow",
				Usage:    "Workflow path",
				Required: true,
				Action: func(c *cli.Context, path cli.Path) error {
					return setting.ParseWorkflow(path)
				},
			},
		},
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
					return wokerStandalone.Execute()
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		return err
	}

	return nil
}
