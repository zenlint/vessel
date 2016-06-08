package setting

import (
	"fmt"
	"io/ioutil"

	"github.com/ghodss/yaml"
)

// GlobalConf global config
type GlobalConf struct {
	AppName     string
	Usage       string
	Version     string
	Author      string
	Email       string
	RuntimePath string
}

// RunTimeConf runtime config
type RunTimeConf struct {
	// Run config
	Run struct {
		runMode string
		logPath string
	}
	// Http config
	HTTP struct {
		ListenMode    string
		HTTPSCertFile string
		HTTPSKeyFile  string
		Host          string
		Port          string
	}
	// Database config
	Database struct {
		Username string
		Password string
		Protocol string
		Host     string
		Port     string
		Schema   string
		Param    map[string]string
	}
	// Etcd config
	Etcd struct {
		Endpoints []map[string]string
		Username  string
		Password  string
	}
	// K8s config
	K8s struct {
		Host string
		Port string
	}
}

var (
	// Global global config
	Global *GlobalConf
	// RunTime runTime config
	RunTime *RunTimeConf
)

// InitConf config init
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
	RunTime = &RunTimeConf{}
	err = yaml.Unmarshal([]byte(runtimeFile), &RunTime)
	if err != nil {
		return err
	}
	fmt.Println(RunTime)
	return nil
}
