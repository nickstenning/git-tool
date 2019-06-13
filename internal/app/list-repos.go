package app

import (
	"fmt"

	"github.com/SierraSoftworks/git-tool/internal/pkg/templates"

	"github.com/SierraSoftworks/git-tool/pkg/filesystem"

	"github.com/urfave/cli"
)

var listReposCommand = cli.Command{
	Name:  "list",
	Usage: "Lists the repositories in your local development environment.",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "quiet,q",
			Usage: "show only the fully qualified repository names",
		},

		cli.BoolFlag{
			Name:  "full",
			Usage: "show all available information about each repository",
		},
	},
	Action: func(c *cli.Context) error {
		mapper := filesystem.Mapper{
			Config: cfg,
		}

		repos, err := mapper.GetRepos()
		if err != nil {
			return err
		}

		for i, repo := range repos {
			if c.Bool("quiet") {
				fmt.Println(templates.RepoQualifiedName(repo, cfg.GetService(repo.Service)))
			} else if c.Bool("full") {
				if i > 0 {
					fmt.Println("---")
				}

				fmt.Println(templates.RepoFullInfo(repo, cfg.GetService(repo.Service)))
			} else {
				fmt.Println(templates.RepoShortInfo(repo, cfg.GetService(repo.Service)))
			}
		}

		return nil
	},
}
