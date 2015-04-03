package cmd

import (
	"github.com/codegangsta/cli"

	api "github.com/containerops/anchor"

	"github.com/containerops/vessel/modules/log"
)

var CmdBuild = cli.Command{
	Name:   "build",
	Usage:  "Build a solution",
	Action: runBuild,
	Flags: []cli.Flag{
		cli.StringFlag{"url", "http://localhost:4000", "API server end point", "VESSEL_URL"},
	},
}

func runBuild(c *cli.Context) {
	client := api.NewClient(c.String("url"))
	opts := api.CreateFlowOptions{
		Name: api.NewString("test"),
	}
	flow, err := client.CreateFlow(opts)
	if err != nil {
		log.Fatal("Fail to create flow: %v", err)
	}
	_ = flow
}
