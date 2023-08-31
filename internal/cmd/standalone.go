package cmd

import (
	"os"

	"github.com/urfave/cli/v2"

	"github.com/galgotech/fermions-workflow/pkg/setting"
	"github.com/galgotech/fermions-workflow/pkg/standalone"
)

func Standalone() error {
	setting := setting.New().(*setting.FermionsSetting)

	app := &cli.App{
		Name:    "fermions-workflow-standalone",
		Usage:   "Fermions workflow runtime standalone",
		Authors: authors,
		Flags:   globalFlags(setting, true, true),
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "Workflow standalone",
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name:     "exec",
						Usage:    "exec workflow",
						Required: false,
						Action: func(c *cli.Context, exec []string) error {
							setting.AddStart(exec)
							return nil
						},
					},
				},
				Action: func(c *cli.Context) error {
					fermionsStandalone, err := standalone.Initialize(setting)
					if err != nil {
						return err
					}
					fermionsStandalone.Execute()
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		return err
	}

	return nil
}
