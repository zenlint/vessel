package kubernetes

import (
	"fmt"

	"k8s.io/kubernetes/pkg/api"
)

// CheckPod check weather the pod spcified by namespace and podname is exist
func CheckPod(namespace string, podName string) bool {

	pods, err := CLIENT.Pods(namespace).List(api.ListOptions{})
	if err != nil {
		fmt.Errorf("List pods err: %v\n", err.Error())
	}

	for _, pod := range pods.Items {
		if pod.Name == podName {
			return true
		}
	}
	return false
}

// GetPodPhase get phase of the resource by namespace and podname, return empty string when no pod find
func GetPodStatus(namespace string, podName string) string {

	pods, err := CLIENT.Pods(namespace).List(api.ListOptions{})
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

// CheckService service have no status, once the service are found, it is with running status
func CheckService(namespace string, serviceName string) bool {

	services, err := CLIENT.Services(namespace).List(api.ListOptions{})
	if err != nil {
		fmt.Errorf("List services err: %v\n", err.Error())
	}

	for _, s := range services.Items {
		if s.Name == serviceName {
			return true
		}
	}
	return false
}

// checkRC rc have no status, once the rc are found, it is with running status
func CheckRC(namespace string, rcName string) bool {
	rcs, err := CLIENT.ReplicationControllers(namespace).List(api.ListOptions{})
	if err != nil {
		fmt.Errorf("List rcs err: %v\n", err.Error())
	}

	for _, rc := range rcs.Items {
		if rc.Name == rcName {
			return true
		}
	}
	return false
}
