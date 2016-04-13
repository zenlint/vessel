package kubernetes

import (
	"fmt"

	"k8s.io/kubernetes/pkg/client/restclient"
	"k8s.io/kubernetes/pkg/client/unversioned"
)

var CLIENT *unversioned.Client

func New(hostIp string) {
	clientConfig := restclient.Config{}
	clientConfig.Host = hostIp
	client, err := unversioned.New(&clientConfig)
	if err != nil {
		fmt.Errorf("New unversioned client err: %v!\n", err.Error())
	}

	CLIENT = client
}
