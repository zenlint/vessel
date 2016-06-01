package etcd

import (
	"testing"
	"log"
	"time"
)

var (
	test_path = "etcd/test"
	test_key = "testStr"
)

func TestCreateClient(t *testing.T) {
	log.Println(clientEtcd())
}

func TestEtcdSetDir(t *testing.T) {
	log.Println(EtcdSetDir(test_path))
	log.Println(EtcdGet(test_path))
}

func TestEtcdGet(t *testing.T) {
	log.Println(EtcdGet(test_path + test_key))
}

func TestEtcdSet(t *testing.T) {
	log.Println(EtcdSet(test_path + test_key, "aaaa"))
	log.Println(EtcdGet(test_path + test_key))
}

func TestEtcdSetTTL(t *testing.T) {
	log.Println(EtcdSetTTL(test_path + test_key, "aaaa", 2))
	<-time.After(time.Second * time.Duration(4))
	log.Println(EtcdGet(test_path + test_key))
}

func TestEtcdSetDirTTL(t *testing.T) {
	log.Println(EtcdSetDirTTL(test_path, 2))
	<-time.After(time.Second * time.Duration(4))
	log.Println(EtcdGet(test_path))
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