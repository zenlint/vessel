package kubernetes

import (
	"time"

	"github.com/containerops/vessel/models"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/watch"
	"log"
)

// CheckPod check weather the pod spcified by namespace and podname is exist
func CheckPod(namespace string, podName string) bool {

	pods, err := models.K8sClient.Pods(namespace).List(api.ListOptions{})
	if err != nil {
		log.Printf("List pods err: %v\n", err.Error())
	}

	for _, pod := range pods.Items {
		if pod.Name == podName {
			return true
		}
	}
	return false
}

func getPodIp(namespace string, rcName string, ipArray *[]string) error {
	// pod, err := models.K8sClient.Pods(namespace).Get(podName)

	opts := api.ListOptions{LabelSelector: labels.Set{"app": rcName}.AsSelector()}
	pods, err := models.K8sClient.Pods(namespace).List(opts)
	if err != nil {
		log.Printf("getPodIp err %v\n", err)
		return err
	}
	for i, pod := range pods.Items {
		(*ipArray)[i] = pod.Status.PodIP
	}

	return nil
	/*if err != nil {
		log.Printf("Get pod %v err: %v\n", podName, err)
		return "", err
	}

	return pod.Status.PodIP, nil*/
}

// GetPodPhase get phase of the resource by namespace and podname, return empty string when no pod find
func GetPodStatus(namespace string, podName string) string {
	pods, err := models.K8sClient.Pods(namespace).List(api.ListOptions{})
	if err != nil {
		log.Printf("List pods err: %v\n", err.Error())
	}

	for _, pod := range pods.Items {
		if pod.Name == podName {
			return string(pod.Status.Phase)
		}
	}
	return ""
}

// WatchPodStatus return status of the operation(specified by checkOp) of the pod, OK, TIMEOUT.
func WatchPodStatus(podNamespace string, labelKey string, labelValue string, timeout int64, checkOp string,sum int,ch chan string){
	log.Printf("Enter WatchPodStatus")
	if checkOp != string(watch.Deleted) && checkOp != string(watch.Added) {
		log.Printf("Params checkOp err, checkOp: %v", checkOp)
	}

	opts := api.ListOptions{LabelSelector: labels.Set{labelKey: labelValue}.AsSelector()}
	w, err := models.K8sClient.Pods(podNamespace).Watch(opts)
	if err != nil {
		ch <- Error
		log.Printf("Get watch interface err")
		return
	}

	t := time.NewTimer(time.Second * time.Duration(timeout))
	for count := 0; count < sum;{
		select {
		case event, ok := <-w.ResultChan():
			if !ok {
				log.Printf("Watch err\n")
				ch <- Error
				// Pod have phase, so we have to wait for the phase change to the right status when added
			}else if string(event.Type) == checkOp {
				log.Println(string(event.Type))
				log.Println(event.Object.(*api.Pod).Status.Phase)
				ch <- OK
			}
			count++
		case <-t.C:
			log.Println("WatchRCStatus timeout")
			ch <- Timeout
		}
	}
}
