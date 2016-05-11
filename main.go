package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"

	"github.com/containerops/vessel/cmd"
	"github.com/containerops/vessel/setting"
)

func main() {
	if err := setting.InitConf("./conf/global.yaml", "./conf/runtime.yaml"); err != nil {
		fmt.Printf("Read config error: %v", err.Error())
		return
	}

	cmdWeb := cmd.GetCmdWeb()

	app := cli.NewApp()

	app.Name = setting.Global.AppName
	app.Usage = setting.Global.Usage
	app.Version = setting.Global.Version
	app.Author = setting.Global.Author
	app.Email = setting.Global.Email
	app.Commands = []cli.Command{
		// cmd.CmdWeb,
		cmdWeb,
		// cmd.CmdDatabase,
	}

	app.Flags = append(app.Flags, []cli.Flag{}...)
	app.Run(os.Args)
}
