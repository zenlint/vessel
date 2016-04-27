package setting

import (
	"fmt"
	"io/ioutil"

	"github.com/ghodss/yaml"
)

type GlobalConf struct {
	AppName     string
	Usage       string
	Version     string
	Author      string
	Email       string
	RuntimePath string
}

type RunTimeConf struct {
	Run struct {
		runMode string
		logPath string
	}
	Http struct {
		ListenMode    string
		HttpsCertFile string
		HttpsKeyFile  string
		Host          string
		Port          string
	}
	Database struct {
		Username string
		Password string
		Protocol string
		Host     string
		Port     string
		Schema   string
		Param    map[string]string
	}
	Etcd struct {
		Endpoints []map[string]string
		Username  string
		Password  string
	}
	K8s struct {
		Host string
		Port string
	}
}

var (
	Global  *GlobalConf
	RunTime *RunTimeConf
)

func InitConf(globalFilePath string, runtimeFilePath string) error {

	globalFile, err := ioutil.ReadFile(globalFilePath)
	if err != nil {
		return err
	}
	Global = &GlobalConf{}
	err = yaml.Unmarshal([]byte(globalFile), &Global)
	if err != nil {
		return err
	}

	runtimeFile, err := ioutil.ReadFile(runtimeFilePath)
	if err != nil {
		return err
	}
	// RunTime := RunTimeConf{}
	RunTime = &RunTimeConf{}
	err = yaml.Unmarshal([]byte(runtimeFile), &RunTime)
	if err != nil {
		return err
	}
	fmt.Println(RunTime)
	return nil
}
