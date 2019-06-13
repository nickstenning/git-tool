package app

import (
	"fmt"
	"strings"

	"github.com/urfave/cli"
)

var listAppsCommand = cli.Command{
	Name:  "apps",
	Usage: "Lists the applications which can be launched with the open command.",
	Action: func(c *cli.Context) error {
		for _, app := range cfg.Applications {
			fmt.Printf("- %s (%s)\n", app.Name, strings.Join(append([]string{app.CommandLine}, app.Arguments...), " "))
		}

		return nil
	},
}
