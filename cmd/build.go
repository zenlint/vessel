package cmd

import (
	"github.com/codegangsta/cli"

	api "github.com/dockercn/anchor"

	"github.com/dockercn/vessel/modules/log"
)

var CmdBuild = cli.Command{
	Name:   "build",
	Usage:  "Build a solution",
	Action: runBuild,
	Flags:  []cli.Flag{},
}

func runBuild(c *cli.Context) {
	client := api.NewClient("http://localhost:4000", "")
	opts := api.CreateFlowOptions{
		Name: api.NewString("test"),
	}
	flow, err := client.CreateFlow(opts)
	if err != nil {
		log.Fatal("Fail to create flow: %v", err)
	}
	_ = flow
}
