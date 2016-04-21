package main

import (
	"os"
	"time"

	"github.com/containerops/vessel/module/etcdpersist/db"
	"github.com/containerops/vessel/module/etcdpersist/persist"

	"github.com/codegangsta/cli"
	"github.com/coreos/go-etcd/etcd"
	"github.com/golang/glog"
)

func main() {
	app := cli.NewApp()
	app.Name = "etcd backup"
	app.Usage = "etcd bakcup"
	app.Commands = []cli.Command{
		{
			Name:    "backup",
			Aliases: []string{"b"},
			Usage:   "backup etcd",
			Action:  backup,
			Flags: []cli.Flag{
				cli.IntFlag{Name: "delay"},
				cli.StringFlag{Name: "etcd-config"},
				cli.StringFlag{Name: "dbname"},
				cli.StringFlag{Name: "dbpwd"},
				cli.StringFlag{Name: "dbip"},
				cli.StringFlag{Name: "dbport"},
				cli.StringFlag{Name: "dbaccount"},
			},
		},
		{
			Name:    "install",
			Aliases: []string{"i"},
			Usage:   "install mysql",
			Action:  install,
			Flags: []cli.Flag{
				cli.IntFlag{Name: "test"},
				cli.StringFlag{Name: "load"}},
		},
	}

	app.Run(os.Args)
}

func createEtcd(configPath string) *etcd.Client {
	etcdClient, err := etcd.NewClientFromFile(configPath)

	if err != nil {
		glog.Fatalln("Can not locate configuration file: `"+configPath+"`. Error: ", err)
	}

	success := etcdClient.SyncCluster()
	if !success {
		glog.Fatalln("cannot sync machines")
	}

	return etcdClient
}
func extractetcd(etcdClient persist.EtcdClient) {
	dataSet := persist.DownloadData([]string{"/"}, true, true, etcdClient)
	for _, v := range dataSet {
		v.InertOrUpdate()
	}
}
func backup(c *cli.Context) {
	if err := db.Connect(c.String("dbport"), c.String("dbip"),
		c.String("dbname"), c.String("dbaccount"), c.String("dbpwd")); err != nil {
		glog.Fatalln(err)
	}
	etcdClient := createEtcd(persist.Conf.EtcdConfigPath)
	for {
		select {
		case <-time.After(time.Second * time.Duration(persist.Conf.Delay)):
			extractetcd(etcdClient)
			glog.Infoln("backup running", time.Now().UTC())
		}
	}
}
func install(c *cli.Context) {
	glog.Infoln("install running")
}
