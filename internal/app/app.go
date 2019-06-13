package app

import (
	"github.com/SierraSoftworks/git-tool/pkg/config"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var cfg = config.Default()

// NewApp creates a new command line application for Git-Tool
func NewApp() *cli.App {
	app := cli.NewApp()

	app.Name = "Git-Tool"
	app.Author = "Benjamin Pannell <benjamin@pannell.dev>"
	app.Copyright = "Copyright © Sierra Softworks 2019"
	app.Usage = "Manage your git repositories"
	app.Version = "0.0.0-dev"

	app.Description = "A tool which helps manage your local git repositories and development folders."

	app.Commands = []cli.Command{
		repoInfoCommand,
		openAppCommand,
		newRepoCommand,
		listReposCommand,
		listAppsCommand,
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "config,c",
			EnvVar: "GITTOOL_CONFIG",
			Usage:  "specify the path to your configuration file",
		},
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "enable verbose logging",
		},
	}

	app.Before = func(c *cli.Context) error {
		if c.GlobalString("config") != "" {
			logrus.WithField("config_path", c.GlobalString("config")).Debug("Loading configuration file")
			cfgResult, err := config.Load(c.GlobalString("config"))
			if err != nil {
				return err
			}

			logrus.WithField("config_path", c.GlobalString("config")).Debug("Loaded configuration file")
			cfg = cfgResult
		}

		if c.GlobalBool("verbose") {
			logrus.SetLevel(logrus.DebugLevel)
		}

		return nil
	}

	return app
}
