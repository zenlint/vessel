package etcd

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/containerops/vessel/models"
	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

var (
	etcdClient client.Client
)

const (
	// ClientErr client not start error
	ClientErr = "ETCD client is not start"
)

func getClient() error {
	if etcdClient == nil {
		etcdClient = models.EtcdClient
	}
	if etcdClient == nil {
		return clientErr()
	}
	return nil
}

func clientErr() error {
	return errors.New(ClientErr)
}

// Get from etcd
func Get(key string) (*client.Response, error) {
	if err := getClient(); err != nil {
		return nil, err
	}
	return client.NewKeysAPI(etcdClient).Get(context.Background(), key, nil)
}

// Set to etcd
func Set(key string, value string, opts *client.SetOptions) error {
	if err := getClient(); err != nil {
		return err
	}
	_, err := client.NewKeysAPI(etcdClient).Set(context.Background(), key, value, opts)
	return err
}

// GetValue (string) from etcd
func GetValue(key string) (string, error) {
	resp, err := Get(key)
	if err != nil {
		return "", err
	}
	return resp.Node.Value, nil
}

// SetValue (string) to etcd
func SetValue(key string, value string) error {
	return Set(key, value, nil)
}

// GetJSON data from etcd
func GetJSON(key string, v interface{}) error {
	jsonStr, err := GetValue(key)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(jsonStr), v)
}

// SetJSON data to etcd
func SetJSON(key string, value interface{}) error {
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return SetValue(key, string(jsonBytes))
}

// SetDir to etcd
func SetDir(key string) error {
	return Set(key, "", &client.SetOptions{Dir: true, PrevExist: client.PrevExist})
}

// GetDir data from etcd
func GetDir(key string) (client.Nodes, error) {
	resp, err := Get(key)
	if err != nil {
		return nil, err
	}
	if !resp.Node.Dir {
		return nil, fmt.Errorf("%v is not dir in ETCD", key)
	}
	return resp.Node.Nodes, nil
}

// SetValueTTL to etcd
func SetValueTTL(key string, value string, timeLife uint64) error {
	return Set(key, value, &client.SetOptions{TTL: time.Duration(timeLife)})
}

// SetDirTTL to etcd
func SetDirTTL(key string, timeLife uint64) error {
	return Set(key, "", &client.SetOptions{TTL: time.Duration(timeLife), Dir: true, PrevExist: client.PrevExist})
}

// Watch on etcd
func Watch(key string) (client.Watcher, error) {
	if etcdClient == nil {
		return nil, clientErr()
	}
	return client.NewKeysAPI(etcdClient).Watcher(key, &client.WatcherOptions{
		Recursive: true,
	}), nil
}
