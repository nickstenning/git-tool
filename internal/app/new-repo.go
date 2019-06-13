package app

import (
	"github.com/urfave/cli"
)

var newRepoCommand = cli.Command{
	Name:  "new",
	Usage: "Creates a new repository with the provided name.",
}
