package etcd

import (
	"fmt"

	"github.com/containerops/vessel/models"
)

const (
	// VesselStagePath vessel stage path in etcd
	VesselStagePath = "containerops/vessel/ns_%v/sn_%v"
)

// SaveStage data to etcd
func SaveStage(stage *models.Stage) error {
	stagePath := fmt.Sprintf(VesselStagePath, stage.Namespace, stage.Name)
	return SetJSON(stagePath+"/data", stage)
}

// GetStage data from etcd
func GetStage(stage *models.Stage) error {
	stagePath := fmt.Sprintf(VesselStagePath, stage.Namespace, stage.Name)
	return GetJSON(stagePath+"/data", stage)
}

// SetStageTTL to etcd
func SetStageTTL(stage *models.Stage, timeLife uint64) error {
	stagePath := fmt.Sprintf(VesselStagePath, stage.Namespace, stage.Name)
	return SetDirTTL(stagePath, timeLife)
}

// SetStageStatus to etcd
func SetStageStatus(stage *models.Stage) error {
	stagePath := fmt.Sprintf(VesselStagePath, stage.Namespace, stage.Name)
	return SetValue(stagePath+"/status", stage.Status)
}

// GetStageStatus from etcd
func GetStageStatus(stage *models.Stage) (string, error) {
	stagePath := fmt.Sprintf(VesselStagePath, stage.Namespace, stage.Name)
	value, err := GetValue(stagePath + "/status")
	if err != nil {
		return "", err
	}
	stage.Status = value
	return value, err
}

// SetStageResult to etcd
func SetStageResult(stageResult *models.StageResult) error {
	stagePath := fmt.Sprintf(VesselStagePath, stageResult.Namespace, stageResult.Name)
	return SetJSON(stagePath+"/result", stageResult)
}

// GetStageResult from etcd
func GetStageResult(stageResult *models.StageResult) error {
	stagePath := fmt.Sprintf(VesselStagePath, stageResult.Namespace, stageResult.Name)
	return GetJSON(stagePath+"/result", stageResult)
}
