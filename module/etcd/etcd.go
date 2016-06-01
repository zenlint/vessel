package etcd

import (
	"encoding/json"
	"errors"
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

func EtcdGetResp(key string) (*client.Response, error) {
	return client.NewKeysAPI(etcdClient).Get(context.Background(), key, nil)
}

func EtcdGet(key string) (string, error) {
	resp, err := EtcdGetResp(key)
	if err != nil {
		return "", err
	}
	return resp.Node.Value, nil
}

func EtcdSet(key string, value string) error {
	_, err := client.NewKeysAPI(etcdClient).Set(context.Background(), key, value, nil)
	return err
}

func EtcdGetJson(key string, v interface{}) (error) {
	jsonStr, err := EtcdGet(key)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(jsonStr), v)
}

func EtcdSetJson(key string, value interface{}) error {
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return EtcdSet(key, string(jsonBytes))
}

func EtcdSetDir(key string) error {
	_, err := client.NewKeysAPI(etcdClient).Set(context.Background(), key, "", &client.SetOptions{Dir: true, PrevExist: client.PrevExist})
	return err
}

func EtcdGetDir(key string) (client.Nodes, error) {
	resp, err := EtcdGetResp(key)
	if err != nil {
		return nil, err
	}
	if !resp.Node.Dir {
		return nil, errors.New(fmt.Sprintf("%v is not dir in ETCD", key))
	}
	return resp.Node.Nodes, nil
}

func EtcdSetTTL(key string, value string, timeLife uint64) error {
	_, err := client.NewKeysAPI(etcdClient).Set(context.Background(), key, value, &client.SetOptions{
		TTL:time.Duration(timeLife),
	})
	return err
}

func EtcdSetDirTTL(key string, timeLife uint64) error {
	_, err := client.NewKeysAPI(etcdClient).Set(context.Background(), key, "", &client.SetOptions{TTL:time.Duration(timeLife), Dir: true, PrevExist: client.PrevExist})
	return err
}

func EtcdWatch(key string) client.Watcher {
	return client.NewKeysAPI(etcdClient).Watcher(key, &client.WatcherOptions{
		Recursive:true,
	})
}