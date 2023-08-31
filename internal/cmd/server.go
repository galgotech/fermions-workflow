package cmd

import (
	"os"

	"github.com/urfave/cli/v2"

	"github.com/galgotech/fermions-workflow/pkg/server"
	"github.com/galgotech/fermions-workflow/pkg/setting"
)

func Server() error {
	setting := setting.New().(*setting.FermionsSetting)

	app := &cli.App{
		Name:    "fermions-workflow-server",
		Usage:   "Workflow server",
		Authors: authors,
		Flags:   globalFlags(setting, true, false),
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "",
				Action: func(c *cli.Context) (err error) {
					fermionsServer, err := server.Initialize(setting)
					if err != nil {
						return err
					}
					return fermionsServer.Execute()
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		return err
	}

	return nil
}
