package kubernetes

import (
	"github.com/containerops/vessel/models"
)

func StartPipeline(pipelineVersion *models.PipelineSpecTemplate) error {
	piplineMetadata := pipelineVersion.MetaData
	if _, err := CLIENT.Namespaces().Get(piplineMetadata.Namespace); err != nil {
		if err := CreateNamespace(pipelineVersion); err != nil {
			return err
		}
	}

	if err := CreateRC(pipelineVersion); err != nil {
		return err
	}

	if err := CreateService(pipelineVersion); err != nil {
		return err
	}

	return nil
}

func DeletePipeline(pipelineVersion *models.PipelineSpecTemplate) error {
	return nil
}

func WatchPipelineStatus(pipelineVersion *models.PipelineSpecTemplate, checkOp string, ch chan string) {
	labelKey := "app"
	pipelineMetadata := pipelineVersion.MetaData
	nsLabelValue := pipelineMetadata.Name
	timeout := pipelineMetadata.TimeoutDuration
	namespace := pipelineMetadata.Namespace

	stageSpecs := pipelineVersion.Spec
	length := len(stageSpecs)
	nsCh := make(chan string)
	rcChs := make([]chan string, length)
	serviceChs := make([]chan string, length)

	go WatchNamespaceStatus(labelKey, nsLabelValue, timeout, checkOp, nsCh)
	for i, stageSpec := range stageSpecs {
		go WatchRCStatus(namespace, labelKey, stageSpec.Name, timeout, checkOp, rcChs[i])
		go WatchServiceStatus(namespace, labelKey, stageSpec.Name, timeout, checkOp, serviceChs[i])
	}

	rcRes := make(chan string)
	serviceRes := make(chan string)
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
	// return
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
