package cmd

import (
	"os"

	"github.com/urfave/cli/v2"

	"github.com/galgotech/fermions-workflow/pkg/server"
	"github.com/galgotech/fermions-workflow/pkg/setting"
)

func Server() error {
	setting := setting.New()
	server, err := server.Initialize(setting)
	if err != nil {
		return err
	}

	app := &cli.App{
		Name:  "Workflow server",
		Usage: "",
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
		},
		Commands: []*cli.Command{
			{
				Name:  "run",
				Usage: "",
				Action: func(c *cli.Context) (err error) {
					return server.Execute()
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		return err
	}

	return nil
}
