package kubernetes

import (
	"fmt"

	"github.com/containerops/vessel/models"
	"github.com/containerops/vessel/utils/timer"
)

// CreateStage from kubernetes
func CreateStage(stage *models.Stage, stageCh chan *models.K8sRes, hourglass *timer.Hourglass) {
	if hourglass.GetLeftNanoseconds() < 0 {
		stageCh <- formatResult(models.ResultTimeout, "Start stage in kubernetes timeout")
		return
	}
	has, err := checkRC(stage)
	if err != nil {
		stageCh <- formatResult(models.ResultFailed, err.Error())
		return
	} else if has {
		stageCh <- formatResult(models.ResultFailed, fmt.Sprintf("Replication controller %v already exist", stage.Name))
		return
	}
	if err := createNamespace(stage); err != nil {
		stageCh <- formatResult(models.ResultFailed, err.Error())
		return
	}

	ch := make(chan *models.K8sRes)
	go watchServiceStatus(stage, models.WatchAdded, hourglass, ch)
	go watchRCStatus(stage, models.WatchAdded, hourglass, ch)
	go watchPodStatus(stage, models.WatchAdded, hourglass, ch)

	if err := createService(stage); err != nil {
		stageCh <- formatResult(models.ResultFailed, err.Error())
		return
	}
	if err := createRC(stage); err != nil {
		stageCh <- formatResult(models.ResultFailed, err.Error())
		return
	}
	var res *models.K8sRes
	for count := 3; count > 0; count-- {
		select {
		case res = <-ch:
			if res.Result != models.ResultSuccess {
				stageCh <- res
				return
			}
		}
	}
	stageCh <- res
}

// DeleteStage from kubernetes
func DeleteStage(stage *models.Stage, stageCh chan *models.K8sRes, hourglass *timer.Hourglass) {
	if hourglass.GetLeftNanoseconds() < 0 {
		stageCh <- formatResult(models.ResultFailed, "Delete stage in kubernetes timeout")
		return
	}
	has, err := checkRC(stage)
	if err != nil {
		stageCh <- formatResult(models.ResultFailed, err.Error())
		return
	} else if !has {
		stageCh <- formatResult(models.ResultFailed, fmt.Sprintf("Replication controller %v not start", stage.Name))
		return
	}

	ch := make(chan *models.K8sRes)
	go watchPodStatus(stage, models.WatchDeleted, hourglass, ch)

	if err := deleteService(stage); err != nil {
		stageCh <- formatResult(models.ResultFailed, err.Error())
		return
	}
	if err := deleteRC(stage); err != nil {
		stageCh <- formatResult(models.ResultFailed, err.Error())
		return
	}
	if err := deletePods(stage); err != nil {
		stageCh <- formatResult(models.ResultFailed, err.Error())
		return
	}

	if count, err := getRCCount(stage); err == nil && count == 0 {
		deleteNamespace(stage)
	}
	stageCh <- ch
}
