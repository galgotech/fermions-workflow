package cmd

import (
	"github.com/galgotech/fermions-workflow/pkg/setting"
	"github.com/urfave/cli/v2"
)

var authors = []*cli.Author{
	{
		Name:  "GalgoTech",
		Email: "andre@galgo.tech",
	},
}

func globalFlags(setting *setting.Setting, conf, workflow bool) []cli.Flag {
	flags := []cli.Flag{}

	if conf {
		flags = append(flags, &cli.PathFlag{
			Name:        "conf",
			Usage:       "",
			DefaultText: "conf/default.ini",
			Required:    false,
			Action: func(c *cli.Context, path cli.Path) error {
				return setting.LoadConfig(path)
			},
		})
	}

	if workflow {
		flags = append(flags, &cli.PathFlag{
			Name:     "workflow",
			Usage:    "Workflow path",
			Required: true,
			Action: func(c *cli.Context, path cli.Path) error {
				return setting.ParseWorkflow(path)
			},
		})
	}

	return flags
}
