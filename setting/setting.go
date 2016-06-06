package setting

/*import (
	"fmt"
	"io/ioutil"

	"github.com/ghodss/yaml"
)*/

// GlobalConf
type GlobalConf struct {
	AppName     string
	Usage       string
	Version     string
	Author      string
	Email       string
	RuntimePath string
}

// RunTimeConf
type RunTimeConf struct {
	// Run
	Run struct {
		runMode string
		logPath string
	}
	// Http
	Http struct {
		ListenMode    string
		HttpsCertFile string
		HttpsKeyFile  string
		Host          string
		Port          string
	}
	// Database
	Database struct {
		Username string
		Password string
		Protocol string
		Host     string
		Port     string
		Schema   string
		Param    map[string]string
	}
	// Etcd
	Etcd struct {
		Endpoints []map[string]string
		Username  string
		Password  string
	}
	// K8s
	K8s struct {
		Host string
		Port string
	}
}

var (
	// Global
	Global  *GlobalConf
	// RunTime
	RunTime *RunTimeConf
)

// InitConf
func InitConf(globalFilePath string, runtimeFilePath string) error {

	/*globalFile, err := ioutil.ReadFile(globalFilePath)
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
	fmt.Println(RunTime)*/
	return nil
}
