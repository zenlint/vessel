package kubernetes

import (
	// "encoding/json"
	"fmt"
	"time"

	"github.com/containerops/vessel/models"
	// "k8s.io/kubernetes/pkg/api"
	// "k8s.io/kubernetes/pkg/util/intstr"
)

func GetPipelineBussinessRes(pipelineVersion *models.PipelineVersion, ch chan bool) {
	namespace := pipelineVersion.GetMetadata().Namespace
	timeout := pipelineVersion.GetMetadata().TimeoutDuration
	// Res := true
	for _, stage := range pipelineVersion.StageSpecs {
		podIp, err := getPodIp(namespace, stage.Name)
		if err != nil {
			ch <- false
			return
		}

		port := stage.Port
		statusCheckLink := stage.StatusCheckUrl
		statusCheckInterval := stage.StatusCheckInterval
		statusCheckCount := stage.StatusCheckCount
		checkUrl := fmt.Sprintf("https://%v:%v%v", podIp, port, statusCheckLink)
		t := time.NewTimer(time.Second * time.Duration(timeout))
		podCh := make(chan bool)
		go getPodBussinessRes(statusCheckLink, statusCheckInterval, statusCheckCount, podCh)
		// for {
		select {
		case podRes := <-ch:
			if podRes == false {
				ch <- false
			}
		case <-t.C:
			// consider timeout as err here
			ch <- false
		}
	}

	ch <- true
}

func getPodBussinessRes(checkUrl string, statusCheckInterval int64, statusCheckCount int, ch chan bool) {
	// request to checkUrl and get time
	// select {}
	for i := 0; i < statusCheckCount; i++ {
		if i == 0 && 0 == requestBsRes(checkUrl) {
			ch <- true
			return
			// bsRes := requestBsRes(checkUrl)
		}

		tick := time.NewTimer(time.Duration(statusCheckInterval) * time.Second)
		<-tick.C
		bsRes := requestBsRes(checkUrl)
		if bsRes == 200 {
			ch <- true
			return
		}
	}
	ch <- false
}

// getBsRes : request to checkUrl, get 200:success, 0, ignore, others, failed
func requestBsRes(checkUrl string) int {
	// read res from checkUrl
	return 200
}
