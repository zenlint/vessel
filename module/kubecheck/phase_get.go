package kubecheck

import (
	"fmt"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/client/restclient"
	"k8s.io/kubernetes/pkg/client/unversioned"
)

func GetPodPhase(hostIp string, podName string) string {
	clientConfig := restclient.Config{}
	clientConfig.Host = hostIp
	client, err := unversioned.New(&clientConfig)
	if err != nil {
		fmt.Errorf("New unversioned client err: %v!\n", err.Error())
	}

	pods, err := client.Pods("").List(api.ListOptions{})
	if err != nil {
		fmt.Errorf("List pods err: %v\n", err.Error())
	}

	for _, pod := range pods.Items {
		if pod.Name == podName {
			return string(pod.Status.Phase)
		}
	}
	return ""
}
