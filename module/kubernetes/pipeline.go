package kubernetes

import (
	// "encoding/json"
	// "fmt"

	"github.com/containerops/vessel/models"
	// "k8s.io/kubernetes/pkg/api"
	// "k8s.io/kubernetes/pkg/util/intstr"
)

func StartPipeline(pipelineVersion *models.PipelineVersion) error {
	piplineMetadata := pipelineVersion.MetaData
	if _, err := CLIENT.Namespaces().Get(piplineMetadata.Namespace); err != nil {
		// fmt.Println("111111111111111")
		if err := CreateNamespace(pipelineVersion); err != nil {
			return err
		}
	}
	// fmt.Println("222222222222222222222222")

	/*if status, err := WatchNamespaceStatus("app", piplineMetadata.Name, 30, Added); err != nil || status != "OK" {
		// if status != "OK" {
		return err
		// }
	}
	*/
	if err := CreateRC(pipelineVersion); err != nil {
		return err
	}

	if err := CreateService(pipelineVersion); err != nil {
		return err
	}
	// CLIENT.Pods(namespace).Get(name).Status.PodIP
	// CLIENT.
	//createrc && createservice
	return nil
}

func DeletePipeline(pipelineVersion *models.PipelineVersion) error {
	return nil
}

func GetPipelinePodsIPort(pipelineVersion *models.PipelineVersion, podIps *[]IpPort) error {
	for _, stage := range pipelineVersion.StageSpecs {
		podIp, err := getPodIp(pipelineVersion.GetMetadata().Namespace, stage.Name)
		if err != nil {
			return err
		}

		(*podIps) = append(*podIps, IpPort{Ip: podIp, Port: stage.Port})
	}
	return nil
}

func WatchPipelineStatus(pipelineVersion *models.PipelineVersion, checkOp string, ch chan string) {
	labelKey := "app"
	pipelineMetadata := pipelineVersion.GetMetadata()
	nsLabelValue := pipelineMetadata.Name
	timeout := pipelineMetadata.TimeoutDuration
	namespace := pipelineMetadata.Namespace

	stageSpecs := pipelineVersion.GetSpec()
	length := len(stageSpecs)
	nsCh := make(chan string)
	rcChs := make([]chan string, length)
	// rcArray := make([]string, length)
	serviceChs := make([]chan string, length)
	// serviceArray := make([]string, length)

	go WatchNamespaceStatus(labelKey, nsLabelValue, timeout, checkOp, nsCh)
	for i, stageSpec := range stageSpecs {
		go WatchRCStatus(namespace, labelKey, stageSpec.Name, timeout, checkOp, rcChs[i])
		go WatchServiceStatus(namespace, labelKey, stageSpec.Name, timeout, checkOp, serviceChs[i])
	}

	// nsRes := make(chan string)
	rcRes := make(chan string)
	serviceRes := make(chan string)
	// go waitNamespace(nsCh, nsRes)
	go wait(length, rcChs, rcRes)
	go wait(length, serviceChs, serviceRes)

	ns := OK
	rc := OK
	service := OK
	for i := 0; i < 3; i++ {
		select {
		case ns = <-nsCh:
			if ns == Error || ns == Timeout {
				ch <- ns
				return
			}
			// temp = nameRes
		case rc = <-rcRes:
			if rc == Error || rc == Timeout {
				ch <- rc
				return
			}
		case service = <-serviceRes:
			if service == Error || service == Timeout {
				ch <- service
				return
			}
		}
	}

	ch <- OK
	return
}

func wait(length int, array []chan string, ch chan string) {
	count := 0
	for i := 0; i < length; i++ {
		res := <-array[i]
		if res == Error || res == Timeout {
			ch <- res
			break
		} else {
			count++
		}
	}
	if count == length-1 {
		ch <- OK
	}
}
