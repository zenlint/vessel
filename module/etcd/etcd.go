package etcd

import (
	"fmt"
	"time"
	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

var (
	etcd_connect_path = "http://%s:%s"
	etcdClient client.Client
)

func CreateClient(settingPoints []map[string]string) (err error) {
	if etcdClient != nil {
		return nil
	}
	endPoints := []string{}
	for _, value := range settingPoints {
		endPoints = append(endPoints, fmt.Sprintf(etcd_connect_path, value["host"], value["port"]))
	}

	cfg := client.Config{
		Endpoints:endPoints,
		Transport:client.DefaultTransport,
		HeaderTimeoutPerRequest:time.Second,
	}
	etcdClient, err = client.New(cfg)
	if err != nil {
		return err
	}
	return nil
}

func EtcdSet(key string, value string) error {
	_, err := client.NewKeysAPI(etcdClient).Set(context.Background(), key, value, nil)
	return err
}

func EtcdSetDir(key string) error{
	_, err := client.NewKeysAPI(etcdClient).Set(context.Background(), key, "", &client.SetOptions{
		Dir:true,
		Refresh:true,
	})
	return err
}

func EtcdSetTTL(key string, value string, timeLife uint64, isDir bool) error {
	_, err := client.NewKeysAPI(etcdClient).Set(context.Background(), key, value, &client.SetOptions{
		TTL:time.Duration(timeLife),
		Dir:isDir,
	})
	return err
}

func EtcdSetDirTTL(key string, timeLife uint64) error {
	_, err := client.NewKeysAPI(etcdClient).Set(context.Background(), key, "", &client.SetOptions{
		TTL:time.Duration(timeLife),
		Dir:true,
	})
	return err
}

func EtcdGet(key string) (*client.Response, error) {
	return client.NewKeysAPI(etcdClient).Get(context.Background(), key, nil)
}

func EtcdWatch(key string) client.Watcher {
	return client.NewKeysAPI(etcdClient).Watcher(key, &client.WatcherOptions{
		Recursive:true,
	})
}