package cmd

import (
	"github.com/codegangsta/cli"
)

var CmdBuild = cli.Command{
	Name:   "build",
	Usage:  "Build a solution",
	Action: runWeb,
	Flags:  []cli.Flag{},
}

func runBuild(c *cli.Context) {

}
