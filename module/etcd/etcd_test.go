package etcd

import (
	"testing"
	"log"
	"fmt"
	"time"

	"github.com/containerops/vessel/models"
	"github.com/coreos/etcd/client"
)

func TestCreateClient(t *testing.T) {
	log.Println(clientEtcd())
}

func clientEtcd() (err error) {
	settingPoints := []map[string]string{
		map[string]string{
			"host": "127.0.0.1",
			"port": "4001",
		},
		map[string]string{
			"host": "localhost",
			"port": "4001",
		},
	}
	if etcdClient != nil {
		return nil
	}
	endPoints := []string{}
	for _, value := range settingPoints {
		endPoints = append(endPoints, fmt.Sprintf(models.ETCD_CONNECT_PATH, value["host"], value["port"]))
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

