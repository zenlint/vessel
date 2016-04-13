package models

import (
	"fmt"
	"time"

	"golang.org/x/net/context"

	"github.com/containerops/vessel/setting"
	"github.com/coreos/etcd/client"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	EtcdClient client.Client
	DbClient   *gorm.DB
)

//Init Database
func InitDatabase() error {

	dbParams := ""
	for key, value := range setting.RunTime.Database.Param {
		if len(dbParams) > 0 {
			dbParams = dbParams + "&" + fmt.Sprintf("%s=%s", key, value)
		} else {
			dbParams = "?" + fmt.Sprintf("%s=%s", key, value)
		}
	}

	connString := fmt.Sprintf(
		"%s:%s@%s(%s:%s)/%s%s",
		setting.RunTime.Database.Username,
		setting.RunTime.Database.Password,
		setting.RunTime.Database.Protocol,
		setting.RunTime.Database.Host,
		setting.RunTime.Database.Port,
		setting.RunTime.Database.Schema,
		dbParams)
	var err error
	DbClient, err = gorm.Open("mysql", connString)
	if err != nil {
		return err
	}

	return nil
}

//Sync Database
func SyncDatabase() error {
	return nil
}

func InitEtcd() error {
	if EtcdClient == nil {
		var etcdEndPoints []string
		for _, value := range setting.RunTime.Etcd.Endpoints {
			etcdEndPoints = append(etcdEndPoints, fmt.Sprintf("http://%s:%s", value["host"], value["port"]))
		}

		cfg := client.Config{
			Endpoints: etcdEndPoints,
			Transport: client.DefaultTransport,
			// set timeout per request to fail fast when the target endpoint is unavailable
			HeaderTimeoutPerRequest: time.Second,
		}
		var err error
		EtcdClient, err = client.New(cfg)
		if err != nil {
			return err
		}
	}
	return nil

}

//Sync ETCD
func SyncETCD() error {
	return nil
}

func EtcdSet(key, value string) error {
	if EtcdClient == nil {
		InitEtcd()
	}
	kapi := client.NewKeysAPI(EtcdClient)
	_, err := kapi.Set(context.Background(), key, value, nil)
	return err
}

func EtcdGet(key string) (*client.Response, error) {
	if EtcdClient == nil {
		InitEtcd()
	}
	kapi := client.NewKeysAPI(EtcdClient)
	return kapi.Get(context.Background(), key, nil)
}

func EtcdWatch(path string) client.Watcher {
	if EtcdClient == nil {
		InitEtcd()
	}
	kapi := client.NewKeysAPI(EtcdClient)
	return kapi.Watcher(path, &client.WatcherOptions{
		Recursive: true,
	})
}
