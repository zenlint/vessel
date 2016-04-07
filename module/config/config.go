package config

import (
	"fmt"

	"github.com/astaxie/beego/config"
)

var (
	conf     config.Configer
	EtcdHost string
	EtcdPort string
)

func Set(path string) error {
	var err error
	conf, err = config.NewConfig("ini", path)
	if err != nil {
		return fmt.Errorf("Read %s error: %v", path, err.Error())
	}

	if etcdHost := conf.String("etcd::host"); etcdHost != "" {
		EtcdHost = etcdHost
	} else {
		err = fmt.Errorf("etcd::host config value is null")
	}

	if etcdPort := conf.String("etcd::port"); etcdPort != "" {
		EtcdPort = etcdPort
	} else {
		err = fmt.Errorf("etcd::port config value is null")
	}

	return err
}

// func Get() (config.Configer, error) {
// 	if conf != nil {
// 		return conf, nil
// 	} else {
// 		return nil, fmt.Errorf("Please Set(path string)")
// 	}
// }
