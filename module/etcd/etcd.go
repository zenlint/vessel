package etcd

import (
	"fmt"
	"time"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/coreos/etcd/client"

	"github.com/containerops/vessel/module/config"
)

/*
	etcd path /vessel/ws-xxx/pj-xxx/pl-xxx/plv-xxx/stage-xxx/

	plv-xxx  -> k8s namespace

	demo:
	/containerops/vessel/ws-xxx/pj-xxx/pl-xxx1/stage/stage1/…
	/containerops/vessel/ws-xxx/pj-xxx/pl-xxx1/stage/stage2…
	/containerops/vessel/ws-xxx/pj-xxx/pl-xxx2/stage/stage1/…
	/containerops/vessel/ws-xxx/pj-xxx/pl-xxx2/stage/stage2…
	/containerops/vessel/ws-xxx/pj-xxx/pl-xxx1/plv-xxx/stagev-xxx/name
	/containerops/vessel/ws-xxx/pj-xxx/pl-xxx1/plv-xxx/stagev-xxx/dependence/Dependence1ServicesName <—need watch
	/containerops/vessel/ws-xxx/pj-xxx/pl-xxx1/plv-xxx/stagev-xxx/dependence/Dependence2ServicesName <—need watch
	/containerops/vessel/ws-xxx/pj-xxx/pl-xxx1/plv-xxx/stagev-xxx/check/check_status_url
	/containerops/vessel/ws-xxx/pj-xxx/pl-xxx1/plv-xxx/stagev-xxx/check/check_status_interval
	/containerops/vessel/ws-xxx/pj-xxx/pl-xxx1/plv-xxx/stagev-xxx/check/check_status_count
*/

var ETCDCLI client.Client

func start() {
	etcdEndPoint := fmt.Sprintf("http://%s:%s", config.EtcdHost, config.EtcdPort)
	cfg := client.Config{
		Endpoints: []string{etcdEndPoint},
		Transport: client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}
	cli, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	ETCDCLI = cli
}

func Set(key, value string) error {
	if ETCDCLI == nil {
		start()
	}
	kapi := client.NewKeysAPI(ETCDCLI)
	_, err := kapi.Set(context.Background(), key, value, nil)
	return err
}

func Get(key string) (*client.Response, error) {
	if ETCDCLI == nil {
		start()
	}
	kapi := client.NewKeysAPI(ETCDCLI)
	return kapi.Get(context.Background(), key, nil)
}

func Watch(path string) client.Watcher {
	if ETCDCLI == nil {
		start()
	}
	kapi := client.NewKeysAPI(ETCDCLI)
	return kapi.Watcher(path, &client.WatcherOptions{
		Recursive: true,
	})
}
