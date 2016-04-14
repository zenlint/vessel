package cmd

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/codegangsta/cli"
	"gopkg.in/macaron.v1"

	"github.com/containerops/vessel/models"
	"github.com/containerops/vessel/setting"
	"github.com/containerops/vessel/web"
)

var CmdWeb = cli.Command{
	Name:        "web",
	Usage:       "start vessel web service",
	Description: "vessel is a CI module.",
	Action:      runWeb,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "address",
			Value: setting.RunTime.Http.Host,
			Usage: "web service listen ip, default is 0.0.0.0; if listen with Unix Socket, the value is sock file path.",
		},
		cli.StringFlag{
			Name:  "port",
			Value: setting.RunTime.Http.Port,
			Usage: "web service listen at port 80; if run with https will be 443.",
		},
	},
}

func runWeb(c *cli.Context) {
	if err := models.InitDatabase(); err != nil {
		fmt.Println(err)
		return
	}
	if err := models.InitEtcd(); err != nil {
		fmt.Println(err)
		return
	}
	m := macaron.New()

	//Set Macaron Web Middleware And Routers
	web.SetVesselMacaron(m)

	switch setting.RunTime.Http.ListenMode {
	case "http":
		listenaddr := fmt.Sprintf("%s:%s", c.String("address"), c.String("port"))
		if err := http.ListenAndServe(listenaddr, m); err != nil {
			fmt.Printf("Start http service error: %v", err.Error())
		}
		break
	case "https":
		listenaddr := fmt.Sprintf("%s:%s", c.String("address"), c.String("port"))
		server := &http.Server{Addr: listenaddr, TLSConfig: &tls.Config{MinVersion: tls.VersionTLS10}, Handler: m}
		if err := server.ListenAndServeTLS(setting.RunTime.Http.HttpsCertFile, setting.RunTime.Http.HttpsKeyFile); err != nil {
			fmt.Printf("Start Dockyard https service error: %v", err.Error())
		}
		break
	default:
		break
	}
}
