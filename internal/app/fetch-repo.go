package app

import (
	"github.com/urfave/cli"
)

var fetchRepoCommand = cli.Command{
	Name:  "fetch",
	Usage: "Fetches the latest version of a remote repository.",
}
