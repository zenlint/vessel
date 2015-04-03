package main

import (
	"os"
	"runtime"

	"github.com/codegangsta/cli"

	"github.com/containerops/vessel/cmd"
	"github.com/containerops/vessel/modules/setting"
)

const APP_VER = "0.0.3.0310"

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	setting.AppVer = APP_VER
	cli.AppHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{.Name}} {{if .Flags}}[global options] {{end}}command{{if .Flags}} [command options]{{end}} [arguments...]

VERSION:
   {{.Version}}

COMMANDS:
   {{range .Commands}}{{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}
   {{end}}{{if .Flags}}
GLOBAL OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}{{end}}
`
}

func main() {
	app := cli.NewApp()
	app.Name = "Vessel"
	app.Usage = "Vessel is a command line utility tool."
	app.Version = APP_VER
	app.Commands = []cli.Command{
		cmd.CmdWeb,
		cmd.CmdBuild,
	}
	app.Flags = append(app.Flags, []cli.Flag{}...)
	app.Run(os.Args)
}

func main1() {
	// log.Println("Creating solution...")
	// sln, err := sln.NewSolutionFromFile(os.Args[1])
	// if err != nil {
	// 	log.Fatalf("Fail to create solution: %v", err)
	// }

	// stage := flow.NewStage()
	// stage.SetJob(sln)
	// if err = stage.Run(); err != nil {
	// 	log.Fatalf("Fail to run stage: %v", err)
	// }

	// log.Println("Starting container...")
	// if err = sln.Start(imageId); err != nil {
	// 	log.Fatalf("Error starting solution: %v", err)
	// }
}
