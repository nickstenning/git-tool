package app

import (
	"fmt"

	"github.com/SierraSoftworks/git-tool/internal/pkg/templates"

	"github.com/SierraSoftworks/git-tool/pkg/filesystem"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

var repoInfoCommand = cli.Command{
	Name:  "info",
	Usage: "Gets the information pertaining to a specific repository.",
	Flags: []cli.Flag{},
	Action: func(c *cli.Context) error {
		service, repo, err := filesystem.GetRepo(cfg, c.Args().First())
		if err != nil {
			return err
		}

		if repo == nil {
			return errors.New("could not find repository")
		}

		fmt.Println(templates.RepoFullInfo(repo, service))

		return nil
	},
}
