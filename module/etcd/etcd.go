package etcd

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
	"github.com/containerops/vessel/models"
)

var (
	etcdClient client.Client
)

const (
	ETCD_CLIENT_ERR = "ETCD client is not start"
)

func getClient() error {
	if etcdClient == nil {
		etcdClient = models.EtcdClient
	}
	if etcdClient == nil {
		return etcdClientErr()
	}
	return nil
}

// EtcdGet Get data from etcd by key
func EtcdGet(key string) (*client.Response, error) {
	if err := getClient(); err != nil {
		return nil, err
	}
	return client.NewKeysAPI(etcdClient).Get(context.Background(), key, nil)
}

// EtcdSet Set data to etcd  by key
func EtcdSet(key string, value string, opts *client.SetOptions) error {
	if err := getClient(); err != nil {
		return err
	}
	_, err := client.NewKeysAPI(etcdClient).Set(context.Background(), key, value, opts)
	return err
}

// EtcdGetValue Get string from etcd by key
func EtcdGetValue(key string) (string, error) {
	resp, err := EtcdGet(key)
	if err != nil {
		return "", err
	}
	return resp.Node.Value, nil
}

// EtcdSetValue Set string to etcd by key
func EtcdSetValue(key string, value string) error {
	return EtcdSet(key, value, nil)
}

// EtcdGetJson Get json data from etcd by key
func EtcdGetJson(key string, v interface{}) (error) {
	jsonStr, err := EtcdGetValue(key)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(jsonStr), v)
}

// EtcdSetJson Set json data to etcd by key
func EtcdSetJson(key string, value interface{}) error {
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return EtcdSetValue(key, string(jsonBytes))
}

// EtcdSetDir Set dir to etcd by key
func EtcdSetDir(key string) error {
	return EtcdSet(key, "", &client.SetOptions{Dir: true, PrevExist: client.PrevExist})
}

// EtcdGetDir Get dir data from etcd by key
func EtcdGetDir(key string) (client.Nodes, error) {
	resp, err := EtcdGet(key)
	if err != nil {
		return nil, err
	}
	if !resp.Node.Dir {
		return nil, errors.New(fmt.Sprintf("%v is not dir in ETCD", key))
	}
	return resp.Node.Nodes, nil
}

// EtcdSetTTL Set data TTL to etcd by key
func EtcdSetTTL(key string, value string, timeLife uint64) error {
	return EtcdSet(key, value, &client.SetOptions{TTL:time.Duration(timeLife)})
}

// EtcdSetDirTTL Set dir TTL to etcd by key
func EtcdSetDirTTL(key string, timeLife uint64) error {
	return EtcdSet(key, "", &client.SetOptions{TTL:time.Duration(timeLife), Dir: true, PrevExist: client.PrevExist})
}

// EtcdWatch Watch from etcd by key
func EtcdWatch(key string) (client.Watcher, error) {
	if etcdClient == nil {
		return nil, etcdClientErr()
	}
	return client.NewKeysAPI(etcdClient).Watcher(key, &client.WatcherOptions{
		Recursive:true,
	}), nil
}

func etcdClientErr() error {
	return errors.New(ETCD_CLIENT_ERR)
}