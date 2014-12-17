package cmd

import (
	"github.com/Scalingo/cli/apps"
	"github.com/codegangsta/cli"
)

var (
	CreateCommand = cli.Command{
		Name:        "create",
		Category:    "Global",
		ShortName:   "c",
		Description: "Create a new app:\n   Example:\n     'scalingo create mynewapp'",
		Usage:       "Create a new app",
		Action: func(c *cli.Context) {
			if len(c.Args()) != 1 {
				cli.ShowCommandHelp(c, "create")
			} else {
				apps.Create(c.Args()[0])
			}
		},
	}
)