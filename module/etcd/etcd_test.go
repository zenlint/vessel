package etcd

import (
	"testing"
	"log"
)

func TestCreateClient(t *testing.T) {
	log.Println(clientEtcd())
}

func clientEtcd() error {
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
	return CreateClient(settingPoints)
}