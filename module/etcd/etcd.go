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

		/vessel/ws-xxx/pj-xxx/pl-xxx/plv-xxx/stage-xxx/name
		/vessel/ws-xxx/pj-xxx/pl-xxx/plv-xxx/stage-xxx/dependence/Dependence1ServicesName <--need watch
		/vessel/ws-xxx/pj-xxx/pl-xxx/plv-xxx/stage-xxx/dependence/Dependence2ServicesName <--need watch

		/vessel/ws-xxx/pj-xxx/pl-xxx/plv-xxx/stage-xxx/status_check_url
		/vessel/ws-xxx/pj-xxx/pl-xxx/plv-xxx/stage-xxx/status_check_interval
		/vessel/ws-xxx/pj-xxx/pl-xxx/plv-xxx/stage-xxx/status_check_count

		/vessel/ws-xxx/pj-xxx/pl-xxx/plv-xxx/stage-xxx/check/check_status_url
		/vessel/ws-xxx/pj-xxx/pl-xxx/plv-xxx/stage-xxx/check/check_status_interval
		/vessel/ws-xxx/pj-xxx/pl-xxx/plv-xxx/stage-xxx/check/check_status_count

		/vessel/ws-xxx/pj-xxx/pl-xxx/plv-xxx/stage-xxx/...
	*/
	log.Error("Setting /vessel/ws-xxx/pj-xxx/pl-xxx/plv-xxx/stage-xxx value")

	// set "/foo" key with "bar" value
	// log.Error("Setting '/foo' key with 'bar' value")
	// log.Print("Setting '/foo' key with 'bar' value")
	resp, err := kapi.Set(context.Background(), "/vessel/ws-xxx/pj-xxx/pl-xxx/plv-xxx/stage-xxx/name", "stage-services-name-xxx", nil)
	if err != nil {
		// log.Fatal(err)
		log.Error(err)

	} else {
		// print common key info
		// log.Printf("Set is done. Metadata is %q\n", resp)
		log.Error("Set is done. Metadata is ", resp)

	}
	// // get "/foo" key's value
	// log.Print("Getting '/foo' key value")
	// resp, err = kapi.Get(context.Background(), "/foo", nil)
	// if err != nil {
	// 	log.Fatal(err)
	// } else {
	// 	// print common key info
	// 	log.Printf("Get is done. Metadata is %q\n", resp)
	// 	// print value
	// 	log.Printf("%q key has %q value\n", resp.Node.Key, resp.Node.Value)
	// }
}
