package kubernetes

import (
	"github.com/containerops/vessel/models"
	"log"
)

func StartPipeline(pipelineVersion *models.PipelineSpecTemplate, stageName string) error {
	piplineMetadata := pipelineVersion.MetaData
	if _, err := models.K8sClient.Namespaces().Get(piplineMetadata.Namespace); err != nil {
		if err := CreateNamespace(pipelineVersion); err != nil {
			return err
		}
	}

	if err := CreateRC(pipelineVersion, stageName); err != nil {
		return err
	}

	if err := CreateService(pipelineVersion, stageName); err != nil {
		return err
	}

	return nil
}

func DeletePipeline(pipelineVersion *models.PipelineSpecTemplate) error {
	meta := pipelineVersion.MetaData
	specs := pipelineVersion.Spec

	for _, spec := range specs {
		models.K8sClient.ReplicationControllers(meta.Namespace).Delete(spec.Name)
		models.K8sClient.Services(meta.Namespace).Delete(spec.Name)
	}

	models.K8sClient.Namespaces().Delete(meta.Namespace)

	return nil
}

func WatchPipelineStatus(pipelineVersion *models.PipelineSpecTemplate, stageName string, checkOp string, ch chan string) {
	log.Println("Enter WatchPipelineStatus")
	labelKey := "app"
	pipelineMetadata := pipelineVersion.MetaData
	// nsLabelValue := pipelineMetadata.Name
	timeout := pipelineMetadata.TimeoutDuration
	namespace := pipelineMetadata.Namespace

	// stageSpecs := pipelineVersion.Spec
	// length := len(stageSpecs)
	// 0423 nsCh := make(chan string)
	//rcCh := make([]chan string, length)
	//serviceCh := make([]chan string, length)
	//0423
	// go WatchNamespaceStatus(labelKey, nsLabelValue, timeout, checkOp, nsCh)
	// rcCh := make(chan string, length)
	// serviceCh := make(chan string, length)

	// for _, stageSpec := range stageSpecs {
	//podSum := pipelineMetadata
	podCh := make(chan string)
	rcCh := make(chan string)
	serviceCh := make(chan string)

	//todo:watch all pod start

	podCount := 0
	for _,specItem := range pipelineVersion.Spec{
		if specItem.Name == stageName{
			podCount = specItem.Replicas
		}
	}
	go WatchPodStatus(namespace, labelKey, stageName, timeout, checkOp, podCount, podCh)
	go WatchRCStatus(namespace, labelKey, stageName, timeout, checkOp, rcCh)
	go WatchServiceStatus(namespace, labelKey, stageName, timeout, checkOp, serviceCh)
	// }

	//rcRes := make(chan string)
	// serviceRes := make(chan string)
	// go wait(length, rcChs, rcRes)
	// go wait(length, serviceChs, serviceRes)

	pod := OK
	rc := OK
	service := OK
	rcCount := 0
	serviceCount := 0
	for i := 0; i < 2 + podCount; i++ {
		select {
		case pod = <-podCh:
			if pod == Error || pod == Timeout {
				log.Println("Get watch pod event err or timeout")
				ch <- pod
				return
			}
			log.Println("Get watch pod event OK")
		case rc = <-rcCh:
			if rc == Error || rc == Timeout {
				log.Println("Get watch rc event err or timeout")
				ch <- rc
				return
			}
			rcCount++
			log.Println("Get watch rc event OK count ", rcCount)
		case service = <-serviceCh:
			if service == Error || service == Timeout {
				log.Println("Get watch service event err or timeout")
				ch <- service
				return
			}
			serviceCount++
			log.Println("Get watch service event ok count ", serviceCount)
		}
	}

	log.Println("WatchPipelineStatus return OK")
	ch <- OK
	// return
}

/*func wait(length int, array []chan string, ch chan string) {
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
*/
