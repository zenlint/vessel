package kubernetes

import (
	"errors"
	"fmt"
	"time"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/watch"
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

func getPodIp(namespace string, podName string) (string, error) {
	pod, err := CLIENT.Pods(namespace).Get(podName)
	if err != nil {
		fmt.Errorf("Get pod %v err: %v\n", podName, err)
		return "", err
	}

	return pod.Status.PodIP, nil
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

// WatchPodStatus return status of the operation(specified by checkOp) of the pod, OK, TIMEOUT.
func WatchPodStatus(podNamespace string, labelKey string, labelValue string, timeout int, checkOp string) (string, error) {
	if checkOp != string(watch.Deleted) && checkOp != string(watch.Added) {
		fmt.Errorf("Params checkOp err, checkOp: %v", checkOp)
	}

	opts := api.ListOptions{LabelSelector: labels.Set{labelKey: labelValue}.AsSelector()}
	w, err := CLIENT.Pods(podNamespace).Watch(opts)
	if err != nil {
		fmt.Errorf("Get watch interface err")
		return "", err
	}

	t := time.NewTimer(time.Second * time.Duration(timeout))
	for {
		select {
		case event, ok := <-w.ResultChan():
			if !ok {
				fmt.Errorf("Watch err\n")
				return "", errors.New("error occours from watch chanle")
			}
			// Pod have phase, so we have to wait for the phase change to the right status when added
			if string(event.Type) == checkOp {
				if (checkOp == string(watch.Deleted)) || ((checkOp != string(watch.Deleted)) &&
					(event.Object.(*api.Pod).Status.Phase == "running")) {
					return "OK", nil
				}
			}
		case <-t.C:
			return "TIMEOUT", nil
		}
	}
}
