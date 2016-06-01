package etcd

import (
	"fmt"

	"github.com/containerops/vessel/models"
)

var (
	vessel_stage_etcd_path = "containerops/vessel/ns_%v/sn_%v"
)

func SaveStage(stage *models.Stage) error {
	stagePath := fmt.Sprintf(vessel_stage_etcd_path, stage.Namespace, stage.Name)
	return EtcdSetJson(stagePath + "/data", stage);
}

func GetStage(stage *models.Stage) error {
	stagePath := fmt.Sprintf(vessel_stage_etcd_path, stage.Namespace, stage.Name)
	return EtcdGetJson(stagePath + "/data", stage)
}

func SetStageTTL(stage *models.Stage, timeLife uint64) error {
	stagePath := fmt.Sprintf(vessel_stage_etcd_path, stage.Namespace, stage.Name)
	return EtcdSetDirTTL(stagePath, timeLife)
}

func SetStageStatus(stage *models.Stage) error {
	stagePath := fmt.Sprintf(vessel_stage_etcd_path, stage.Namespace, stage.Name)
	return EtcdSet(stagePath + "/status", stage.Status)
}

func GetStageStatus(stage *models.Stage) (string, error) {
	stagePath := fmt.Sprintf(vessel_stage_etcd_path, stage.Namespace, stage.Name)
	value, err := EtcdGet(stagePath + "/status")
	if err != nil {
		return "", err
	}
	stage.Status = value
	return value, err
}

func SetStageResult(stageResult *models.StageResult) error {
	stagePath := fmt.Sprintf(vessel_stage_etcd_path, stageResult.Namespace, stageResult.Name)
	return EtcdSetJson(stagePath + "/result", stageResult)
}

func GetStageResult(stageResult *models.StageResult) error {
	stagePath := fmt.Sprintf(vessel_stage_etcd_path, stageResult.Namespace, stageResult.Name)
	return EtcdGetJson(stagePath + "/result", stageResult)
}