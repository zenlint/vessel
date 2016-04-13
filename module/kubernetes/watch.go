package kubernetes

import (
	"errors"
	"fmt"
	"time"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/watch"
)

/*
Example:
package main

import (
	"fmt"
	"github.com/zenlinTechnofreak/vessel/module/kube"
)

func main() {
	kube.New("127.0.0.1:8080")
	b, err := kube.WatchPodStatus("", "app", "nginx", 30, "ADDED")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(b)
}

*/

// WatchPodStatus return status of the operation(specified by checkOp) of the pod, OK, TIMEOUT.
func WatchPodStatus(podNamespace string, labelKey string, labelValue string, timeout int, checkOp string) (string, error) {
	if checkOp != string(watch.Deleted) && checkOp != string(watch.Added) {
		fmt.Errorf("Params checkOp err, checkOp: %v", checkOp)
	}

	//opts := api.ListOptions{FieldSelector: fields.Set{"kind": "pod"}.AsSelector()}
	opts := api.ListOptions{LabelSelector: labels.Set{labelKey: labelValue}.AsSelector()}

	w, err := CLIENT.Pods(podNamespace).Watch(opts)
	if err != nil {
		fmt.Errorf("Get watch interface err")
	}

	t := time.NewTimer(time.Second * time.Duration(timeout))

	for {
		select {
		case event, ok := <-w.ResultChan():
			//fmt.Println(event.Type)
			if !ok {
				fmt.Errorf("Watch err\n")
				return "", errors.New("error occours from watch chanle")
			}
			//fmt.Println(event.Type)
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

// WatchServiceStatus return status of the operation(specified by checkOp) of the pod, OK, TIMEOUT.
func WatchServiceStatus(Namespace string, labelKey string, labelValue string, timeout int, checkOp string) (string, error) {
	if checkOp != string(watch.Deleted) && checkOp != string(watch.Added) {
		fmt.Errorf("Params checkOp err, checkOp: %v", checkOp)
	}

	//opts := api.ListOptions{FieldSelector: fields.Set{"kind": "pod"}.AsSelector()}
	opts := api.ListOptions{LabelSelector: labels.Set{labelKey: labelValue}.AsSelector()}

	w, err := CLIENT.Services(Namespace).Watch(opts)
	if err != nil {
		fmt.Errorf("Get watch interface err")
	}

	t := time.NewTimer(time.Second * time.Duration(timeout))

	for {
		select {
		case event, ok := <-w.ResultChan():
			//fmt.Println(event.Type)
			if !ok {
				fmt.Errorf("Watch err\n")
				return "", errors.New("error occours from watch chanle")
			}
			//fmt.Println(event.Type)
			if string(event.Type) == checkOp {
				return "OK", nil
			}

		case <-t.C:
			return "TIMEOUT", nil
		}
	}
}

// WatchServiceStatus return status of the operation(specified by checkOp) of the pod, OK, TIMEOUT.
func WatchRCStatus(Namespace string, labelKey string, labelValue string, timeout int, checkOp string) (string, error) {
	if checkOp != string(watch.Deleted) && checkOp != string(watch.Added) {
		fmt.Errorf("Params checkOp err, checkOp: %v", checkOp)
	}

	//opts := api.ListOptions{FieldSelector: fields.Set{"kind": "pod"}.AsSelector()}
	opts := api.ListOptions{LabelSelector: labels.Set{labelKey: labelValue}.AsSelector()}

	w, err := CLIENT.ReplicationControllers(Namespace).Watch(opts)
	if err != nil {
		fmt.Errorf("Get watch interface err")
	}

	t := time.NewTimer(time.Second * time.Duration(timeout))

	for {
		select {
		case event, ok := <-w.ResultChan():
			//fmt.Println(event.Type)
			if !ok {
				fmt.Errorf("Watch err\n")
				return "", errors.New("error occours from watch chanle")
			}
			//fmt.Println(event.Type)
			if string(event.Type) == checkOp {
				return "OK", nil
			}

		case <-t.C:
			return "TIMEOUT", nil
		}
	}
}

/*func WatchPod(podName string, podNamespace string, c chan string) {

	//opts := api.ListOptions{FieldSelector: fields.Set{"kind": "pod"}.AsSelector()}
	opts := api.ListOptions{LabelSelector: labels.Set{"app": "nginx"}.AsSelector()}

	w, err := CLIENT.Pods(podNamespace).Watch(opts)
	if err != nil {
		fmt.Errorf("Get watch interface err")
	}

	for {
		event, ok := <-w.ResultChan()

		if !ok {
			// Resource was deleted, and chanle closed, so return to main programme
			return
		}
		switch event.Type {
		case "DELETED":
			c <- DELETED
			w.Stop()
		case "ERROR":
			c <- ERROR
			w.Stop()
		default:
			if event.Object.(*api.Pod).Status.Phase == "running" {
				c <- RUNNING
			} else {
				c <- PENDING
			}
		}
	}
}
*/
