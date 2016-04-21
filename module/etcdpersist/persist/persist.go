package persist

import (
	"github.com/containerops/vessel/module/etcdpersist/db"

	"github.com/coreos/go-etcd/etcd"
	"github.com/golang/glog"
)

type EtcdClient interface {
	Get(key string, sort, recursive bool) (*etcd.Response, error)
	Set(key string, value string, ttl uint64) (*etcd.Response, error)
	SetDir(key string, ttl uint64) (*etcd.Response, error)
}

// load data from etcd cache
func DownloadData(keys []string, sorted bool, recursive bool, etcdClient EtcdClient) []*db.Tb_etcd_backup {
	persistData := make([]*db.Tb_etcd_backup, 0)
	for _, key := range keys {
		response, err := etcdClient.Get(key, sorted, recursive)
		if err != nil {
			glog.Infoln("Trying to get the following key: "+key+". Error: ", err)
		}
		persistData = append(persistData, extractNodes(response.Node, recursive)...)
	}
	return persistData
}

func extractNodes(node *etcd.Node, recursive bool) []*db.Tb_etcd_backup {
	backupKeys := make([]*db.Tb_etcd_backup, 0)
	if recursive {
		backupKeys = backupNodes(node)
	} else {
		backupKeys = append(backupKeys, backupNode(node))
	}
	return backupKeys
}

// copy data to database struct
func backupNode(node *etcd.Node) *db.Tb_etcd_backup {
	key := &db.Tb_etcd_backup{
		Key:            node.Key,
		Modified_index: node.ModifiedIndex,
		Created_index:  node.CreatedIndex,
	}
	if node.Dir {
		key.Dir = 1
	} else {
		key.Dir = 0
	}
	if node.Expiration != nil {
		key.Ttl = node.Expiration.Unix()
	}
	if node.Dir != true && node.Key != "" {
		key.Value = node.Value
	}
	return key
}

func backupNodes(node *etcd.Node) []*db.Tb_etcd_backup {
	backupKeys := []*db.Tb_etcd_backup{}

	if len(node.Nodes) > 0 {
		for _, nodeChild := range node.Nodes {
			backupKeys = append(backupKeys, backupNodes(nodeChild)...)
		}
	} else {
		backupKey := backupNode(node)
		if backupKey.Key != "" {
			backupKeys = append(backupKeys, backupKey)
		}
	}
	return backupKeys
}
