package etcd

import (
	"fmt"
	"time"
	"strings"

	"github.com/containerops/vessel/models"
)

const (
	VESSEL_PIPELINE_ETCD_PATH = "containerops/vessel/ns_%v/pl_%v"
)

// SavePipeline Save pipeline to etcd
func SavePipeline(pipeline *models.Pipeline) error {
	pipelinePath := fmt.Sprintf(VESSEL_PIPELINE_ETCD_PATH, pipeline.Namespace, pipeline.Name)
	if err := EtcdSetJson(pipelinePath + "/data", pipeline); err != nil{
		return err
	}
	return SavePipelineStages(pipeline)
}

// GetPipeline Get pipeline from etcd
func GetPipeline(pipeline *models.Pipeline) error {
	pipelinePath := fmt.Sprintf(VESSEL_PIPELINE_ETCD_PATH, pipeline.Namespace, pipeline.Name)
	if err := EtcdGetJson(pipelinePath + "/data", pipeline); err != nil{
		return err
	}
	return GetPipelineStages(pipeline)
}

// SavePipelineStages Save pipeline stages to etcd
func SavePipelineStages(pipeline *models.Pipeline) error {
	pipelinePath := fmt.Sprintf(VESSEL_PIPELINE_ETCD_PATH, pipeline.Namespace, pipeline.Name)
	return EtcdSetValue(pipelinePath + "/stages", strings.Join(pipeline.Stages, ","))
}

// GetPipelineStages Get pipeline stages from etcd
func GetPipelineStages(pipeline *models.Pipeline) error {
	pipelinePath := fmt.Sprintf(VESSEL_PIPELINE_ETCD_PATH, pipeline.Namespace, pipeline.Name)
	value, err := EtcdGetValue(pipelinePath + "/stages")
	if err != nil {
		return err
	}
	pipeline.Stages = strings.Split(value, ",")
	return nil
}

// SetPipelineTTL Set pipeline TTL to etcd
func SetPipelineTTL(pipeline *models.Pipeline, timeLife uint64) error {
	pipelinePath := fmt.Sprintf(VESSEL_PIPELINE_ETCD_PATH, pipeline.Namespace, pipeline.Name)
	return EtcdSetDirTTL(pipelinePath, timeLife)
}

// SetPipelineStatus Set pipeline status to etcd
func SetPipelineStatus(pipeline *models.Pipeline) error {
	pipelinePath := fmt.Sprintf(VESSEL_PIPELINE_ETCD_PATH, pipeline.Namespace, pipeline.Name)
	return EtcdSetValue(pipelinePath + "/status", pipeline.Status)
}

// GetPipelineStatus Get pipeline status from etcd
func GetPipelineStatus(pipeline *models.Pipeline) (string, error) {
	pipelinePath := fmt.Sprintf(VESSEL_PIPELINE_ETCD_PATH, pipeline.Namespace, pipeline.Name)
	value, err := EtcdGetValue(pipelinePath + "/status")
	if err != nil {
		return "", err
	}
	pipeline.Status = value
	return value, err
}

// SetCreationTimestamp Set pipeline creation timestamp to etcd
func SetCreationTimestamp(pipeline *models.Pipeline) error {
	pipelinePath := fmt.Sprintf(VESSEL_PIPELINE_ETCD_PATH, pipeline.Namespace, pipeline.Name)
	pipeline.CreationTimestamp = time.Now().Format("2006-01-02 15:04:05")
	return EtcdSetValue(pipelinePath + "/creationTimestamp", pipeline.CreationTimestamp)
}

// GetCreationTimestamp Get pipeline creation timestamp from etcd
func GetCreationTimestamp(pipeline *models.Pipeline) (string, error) {
	pipelinePath := fmt.Sprintf(VESSEL_PIPELINE_ETCD_PATH, pipeline.Namespace, pipeline.Name)
	value, err := EtcdGetValue(pipelinePath + "/creationTimestamp")
	if err != nil {
		return "", err
	}
	pipeline.CreationTimestamp = value
	return value, err
}

// SetDeletionTimestamp Set pipeline deletion timestamp to etcd
func SetDeletionTimestamp(pipeline *models.Pipeline) error {
	pipelinePath := fmt.Sprintf(VESSEL_PIPELINE_ETCD_PATH, pipeline.Namespace, pipeline.Name)
	pipeline.DeletionTimestamp = time.Now().Format("2006-01-02 15:04:05")
	return EtcdSetValue(pipelinePath + "/deletiontimestamp", pipeline.DeletionTimestamp)
}

// GetDeletionTimestamp Get pipeline deletion timestamp from etcd
func GetDeletionTimestamp(pipeline *models.Pipeline) (string, error) {
	pipelinePath := fmt.Sprintf(VESSEL_PIPELINE_ETCD_PATH, pipeline.Namespace, pipeline.Name)
	value, err := EtcdGetValue(pipelinePath + "/deletiontimestamp")
	if err != nil {
		return "", err
	}
	pipeline.DeletionTimestamp = value
	return value, err
}

// SetPipelineResult Set pipeline result to etcd
func SetPipelineResult(pipelineResult *models.PipelineResult) error {
	pipelinePath := fmt.Sprintf(VESSEL_PIPELINE_ETCD_PATH, pipelineResult.Namespace, pipelineResult.Name)
	return EtcdSetJson(pipelinePath + "/result", pipelineResult)
}

// GetPipelineResult Get pipeline result from etcd
func GetPipelineResult(pipelineResult *models.PipelineResult) error {
	pipelinePath := fmt.Sprintf(VESSEL_PIPELINE_ETCD_PATH, pipelineResult.Namespace, pipelineResult.Name)
	return EtcdGetJson(pipelinePath + "/result", pipelineResult)
}