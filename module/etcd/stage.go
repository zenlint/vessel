package etcd

import (
	"fmt"

	"github.com/containerops/vessel/models"
)

const (
	VESSEL_STAGE_ETCD_PATH = "containerops/vessel/ns_%v/sn_%v"
)

// SaveStage Save stage to etcd
func SaveStage(stage *models.Stage) error {
	stagePath := fmt.Sprintf(VESSEL_STAGE_ETCD_PATH, stage.Namespace, stage.Name)
	return EtcdSetJson(stagePath + "/data", stage);
}

// GetStage Get stage from etcd
func GetStage(stage *models.Stage) error {
	stagePath := fmt.Sprintf(VESSEL_STAGE_ETCD_PATH, stage.Namespace, stage.Name)
	return EtcdGetJson(stagePath + "/data", stage)
}

// SetStageTTL Set stage TTL to etcd
func SetStageTTL(stage *models.Stage, timeLife uint64) error {
	stagePath := fmt.Sprintf(VESSEL_STAGE_ETCD_PATH, stage.Namespace, stage.Name)
	return EtcdSetDirTTL(stagePath, timeLife)
}

// SetStageStatus Set stage status to etcd
func SetStageStatus(stage *models.Stage) error {
	stagePath := fmt.Sprintf(VESSEL_STAGE_ETCD_PATH, stage.Namespace, stage.Name)
	return EtcdSetValue(stagePath + "/status", stage.Status)
}

// GetStageStatus Get stage status from etcd
func GetStageStatus(stage *models.Stage) (string, error) {
	stagePath := fmt.Sprintf(VESSEL_STAGE_ETCD_PATH, stage.Namespace, stage.Name)
	value, err := EtcdGetValue(stagePath + "/status")
	if err != nil {
		return "", err
	}
	stage.Status = value
	return value, err
}

// SetStageResult Set stage result to etcd
func SetStageResult(stageResult *models.StageResult) error {
	stagePath := fmt.Sprintf(VESSEL_STAGE_ETCD_PATH, stageResult.Namespace, stageResult.Name)
	return EtcdSetJson(stagePath + "/result", stageResult)
}

// GetStageResult Get stage result from etcd
func GetStageResult(stageResult *models.StageResult) error {
	stagePath := fmt.Sprintf(VESSEL_STAGE_ETCD_PATH, stageResult.Namespace, stageResult.Name)
	return EtcdGetJson(stagePath + "/result", stageResult)
}