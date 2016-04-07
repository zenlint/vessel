package etcd

import (
	"fmt"
	"time"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
	"github.com/coreos/etcd/client"

	"github.com/containerops/vessel/module/config"
)

func Save() {
	etcdEndPoint := fmt.Sprintf("http://%s:%s", config.EtcdHost, config.EtcdPort)
	cfg := client.Config{
		Endpoints: []string{etcdEndPoint},
		Transport: client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: time.Second,
	}
	c, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	kapi := client.NewKeysAPI(c)
	/*
		etcd path /vessel/ws-xxx/pj-xxx/pl-xxx/plv-xxx/stage-xxx/

		plv-xxx  -> k8s namespace

		demo:

		/vessel/ws-xxx/pj-xxx/pl-xxx1/stage/stage1/...
		/vessel/ws-xxx/pj-xxx/pl-xxx1/stage/stage2...
		/vessel/ws-xxx/pj-xxx/pl-xxx2/stage/stage1/...
		/vessel/ws-xxx/pj-xxx/pl-xxx2/stage/stage2...

		/vessel/ws-xxx/pj-xxx/pl-xxx1/plv-xxx/stagev-xxx/name
		/vessel/ws-xxx/pj-xxx/pl-xxx1/plv-xxx/stagev-xxx/dependence/Dependence1ServicesName <--need watch
		/vessel/ws-xxx/pj-xxx/pl-xxx1/plv-xxx/stagev-xxx/dependence/Dependence2ServicesName <--need watch

		/vessel/ws-xxx/pj-xxx/pl-xxx1/plv-xxx/stagev-xxx/check/check_status_url
		/vessel/ws-xxx/pj-xxx/pl-xxx1/plv-xxx/stagev-xxx/check/check_status_interval
		/vessel/ws-xxx/pj-xxx/pl-xxx1/plv-xxx/stagev-xxx/check/check_status_count

	*/

	log.Error("Setting /vessel/ws-xxx/pj-xxx/pl-xxx/plv-xxx/stage-xxx value")

	resp, err := kapi.Set(context.Background(), "/vessel/ws-xxx/pj-xxx/pl-xxx/plv-xxx/stage-xxx/name", "stage-services-name-xxx", nil)
	if err != nil {
		log.Error(err)
	} else {
		log.Error("Set is done. Metadata is ", resp)
	}
}
