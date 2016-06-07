package etcd

import (
	"fmt"
	"strings"
	"time"

	"github.com/containerops/vessel/models"
)

const (
	// VesselPipelinePath vessel pipeline path in etcd
	VesselPipelinePath = "containerops/vessel/ns_%v/pl_%v"
)

// SavePipeline data to etcd
func SavePipeline(pipeline *models.Pipeline) error {
	pipelinePath := fmt.Sprintf(VesselPipelinePath, pipeline.Namespace, pipeline.Name)
	if err := SetJSON(pipelinePath+"/data", pipeline); err != nil {
		return err
	}
	return SavePipelineStages(pipeline)
}

// GetPipeline data from etcd
func GetPipeline(pipeline *models.Pipeline) error {
	pipelinePath := fmt.Sprintf(VesselPipelinePath, pipeline.Namespace, pipeline.Name)
	if err := GetJSON(pipelinePath+"/data", pipeline); err != nil {
		return err
	}
	return GetPipelineStages(pipeline)
}

// SavePipelineStages to etcd
func SavePipelineStages(pipeline *models.Pipeline) error {
	pipelinePath := fmt.Sprintf(VesselPipelinePath, pipeline.Namespace, pipeline.Name)
	return SetValue(pipelinePath+"/stages", strings.Join(pipeline.Stages, ","))
}

// GetPipelineStages from etcd
func GetPipelineStages(pipeline *models.Pipeline) error {
	pipelinePath := fmt.Sprintf(VesselPipelinePath, pipeline.Namespace, pipeline.Name)
	value, err := GetValue(pipelinePath + "/stages")
	if err != nil {
		return err
	}
	pipeline.Stages = strings.Split(value, ",")
	return nil
}

// SetPipelineTTL to etcd
func SetPipelineTTL(pipeline *models.Pipeline, timeLife uint64) error {
	pipelinePath := fmt.Sprintf(VesselPipelinePath, pipeline.Namespace, pipeline.Name)
	return SetDirTTL(pipelinePath, timeLife)
}

// SetPipelineStatus to etcd
func SetPipelineStatus(pipeline *models.Pipeline) error {
	pipelinePath := fmt.Sprintf(VesselPipelinePath, pipeline.Namespace, pipeline.Name)
	return SetValue(pipelinePath+"/status", pipeline.Status)
}

// GetPipelineStatus from etcd
func GetPipelineStatus(pipeline *models.Pipeline) (string, error) {
	pipelinePath := fmt.Sprintf(VesselPipelinePath, pipeline.Namespace, pipeline.Name)
	value, err := GetValue(pipelinePath + "/status")
	if err != nil {
		return "", err
	}
	pipeline.Status = value
	return value, err
}

// SetCreationTimestamp to etcd
func SetCreationTimestamp(pipeline *models.Pipeline) error {
	pipelinePath := fmt.Sprintf(VesselPipelinePath, pipeline.Namespace, pipeline.Name)
	pipeline.CreationTimestamp = time.Now().Format("2006-01-02 15:04:05")
	return SetValue(pipelinePath+"/creationTimestamp", pipeline.CreationTimestamp)
}

// GetCreationTimestamp from etcd
func GetCreationTimestamp(pipeline *models.Pipeline) (string, error) {
	pipelinePath := fmt.Sprintf(VesselPipelinePath, pipeline.Namespace, pipeline.Name)
	value, err := GetValue(pipelinePath + "/creationTimestamp")
	if err != nil {
		return "", err
	}
	pipeline.CreationTimestamp = value
	return value, err
}

// SetDeletionTimestamp to etcd
func SetDeletionTimestamp(pipeline *models.Pipeline) error {
	pipelinePath := fmt.Sprintf(VesselPipelinePath, pipeline.Namespace, pipeline.Name)
	pipeline.DeletionTimestamp = time.Now().Format("2006-01-02 15:04:05")
	return SetValue(pipelinePath+"/deletiontimestamp", pipeline.DeletionTimestamp)
}

// GetDeletionTimestamp from etcd
func GetDeletionTimestamp(pipeline *models.Pipeline) (string, error) {
	pipelinePath := fmt.Sprintf(VesselPipelinePath, pipeline.Namespace, pipeline.Name)
	value, err := GetValue(pipelinePath + "/deletiontimestamp")
	if err != nil {
		return "", err
	}
	pipeline.DeletionTimestamp = value
	return value, err
}

// SetPipelineResult to etcd
func SetPipelineResult(pipelineResult *models.PipelineResult) error {
	pipelinePath := fmt.Sprintf(VesselPipelinePath, pipelineResult.Namespace, pipelineResult.Name)
	return SetJSON(pipelinePath+"/result", pipelineResult)
}

// GetPipelineResult from etcd
func GetPipelineResult(pipelineResult *models.PipelineResult) error {
	pipelinePath := fmt.Sprintf(VesselPipelinePath, pipelineResult.Namespace, pipelineResult.Name)
	return GetJSON(pipelinePath+"/result", pipelineResult)
}
