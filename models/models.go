package models

import (
	"fmt"
	"time"

	"github.com/containerops/vessel/setting"
	"k8s.io/kubernetes/pkg/client/restclient"
	"github.com/coreos/etcd/client"
	"k8s.io/kubernetes/pkg/client/unversioned"
)

var (
	EtcdClient client.Client
	K8sClient  *unversioned.Client
)

const (
	ETCD_CONNECT_PATH = "http://%s:%s"
	K8S_CONNECT_PATH = "%v:%v"
)

// InitEtcd
func InitEtcd() error {
	if EtcdClient == nil {
		var etcdEndPoints []string
		for _, value := range setting.RunTime.Etcd.Endpoints {
			etcdEndPoints = append(etcdEndPoints, fmt.Sprintf(ETCD_CONNECT_PATH, value["host"], value["port"]))
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

// InitK8S
func InitK8S() error {
	if K8sClient == nil {
		clientConfig := restclient.Config{}
		host := fmt.Sprintf(K8S_CONNECT_PATH, setting.RunTime.K8s.Host, setting.RunTime.K8s.Port)
		fmt.Println(host)
		clientConfig.Host = host
		// clientConfig.Host = setting.RunTime.Database.Host
		client, err := unversioned.New(&clientConfig)
		if err != nil {
			return err
			fmt.Errorf("New unversioned client err: %v!\n", err.Error())
		}
		K8sClient = client
	}
	return nil
}