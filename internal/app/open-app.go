package app

import (
	"os"

	"github.com/SierraSoftworks/git-tool/pkg/filesystem"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var openAppCommand = cli.Command{
	Name:      "open",
	Usage:     "Opens the requested repository in a specific command.",
	ArgsUsage: "[APP] [REPO]",
	Flags:     []cli.Flag{},
	Action: func(c *cli.Context) error {
		if len(cfg.Applications) == 0 {
			return errors.New("no apps defined in your config")
		}

		app := cfg.GetApp(c.Args().Get(0))
		if app == nil {
			return errors.Errorf("no app called %s in your config", c.Args().Get(0))
		}

		logrus.WithField("app", app.Name).Debug("Found matching app configuration")

		service, repo, err := filesystem.GetRepo(cfg, c.Args().Get(1))
		if err != nil {
			return err
		}

		if repo == nil {
			return errors.New("could not find repository")
		}

		cmd, err := app.GetCmd(service, repo)
		if err != nil {
			return err
		}

		cmd.Stdin = os.Stdin
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout

		return cmd.Run()
	},
}
