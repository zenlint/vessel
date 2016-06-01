package etcd

import (
	"fmt"
	"time"

	"github.com/containerops/vessel/models"
)

var (
	vessel_pipeline_etcd_path = "containerops/vessel/ns_%v/pl_%v"
)

func SavePipeline(pipeline *models.Pipeline) error {
	pipelinePath := fmt.Sprintf(vessel_pipeline_etcd_path, pipeline.Namespace, pipeline.Name)
	return EtcdSetJson(pipelinePath + "/data", pipeline)
}

func GetPipeline(pipeline *models.Pipeline) error {
	pipelinePath := fmt.Sprintf(vessel_pipeline_etcd_path, pipeline.Namespace, pipeline.Name)
	return EtcdGetJson(pipelinePath + "/data", pipeline)
}

func SetPipelineTTL(pipeline *models.Pipeline, timeLife uint64) error {
	pipelinePath := fmt.Sprintf(vessel_pipeline_etcd_path, pipeline.Namespace, pipeline.Name)
	return EtcdSetDirTTL(pipelinePath, timeLife)
}

func SetPipelineStatus(pipeline *models.Pipeline) error {
	pipelinePath := fmt.Sprintf(vessel_pipeline_etcd_path, pipeline.Namespace, pipeline.Name)
	return EtcdSet(pipelinePath + "/status", pipeline.Status)
}

func GetPipelineStatus(pipeline *models.Pipeline) (string, error) {
	pipelinePath := fmt.Sprintf(vessel_pipeline_etcd_path, pipeline.Namespace, pipeline.Name)
	value, err := EtcdGet(pipelinePath + "/status")
	if err != nil {
		return "", err
	}
	pipeline.Status = value
	return value, err
}

func SetCreationTimestamp(pipeline *models.Pipeline) error {
	pipelinePath := fmt.Sprintf(vessel_pipeline_etcd_path, pipeline.Namespace, pipeline.Name)
	pipeline.CreationTimestamp = time.Now().Format("2006-01-02 15:04:05")
	return EtcdSet(pipelinePath + "/creationTimestamp", pipeline.CreationTimestamp)
}

func GetCreationTimestamp(pipeline *models.Pipeline) (string, error) {
	pipelinePath := fmt.Sprintf(vessel_pipeline_etcd_path, pipeline.Namespace, pipeline.Name)
	value, err := EtcdGet(pipelinePath + "/creationTimestamp")
	if err != nil {
		return "", err
	}
	pipeline.CreationTimestamp = value
	return value, err
}

func SetDeletionTimestamp(pipeline *models.Pipeline) error {
	pipelinePath := fmt.Sprintf(vessel_pipeline_etcd_path, pipeline.Namespace, pipeline.Name)
	pipeline.DeletionTimestamp = time.Now().Format("2006-01-02 15:04:05")
	return EtcdSet(pipelinePath + "/deletiontimestamp", pipeline.DeletionTimestamp)
}

func GetDeletionTimestamp(pipeline *models.Pipeline) (string, error) {
	pipelinePath := fmt.Sprintf(vessel_pipeline_etcd_path, pipeline.Namespace, pipeline.Name)
	value, err := EtcdGet(pipelinePath + "/deletiontimestamp")
	if err != nil {
		return "", err
	}
	pipeline.DeletionTimestamp = value
	return value, err
}

func SetPipelineResult(pipelineResult *models.PipelineResult) error {
	pipelinePath := fmt.Sprintf(vessel_pipeline_etcd_path, pipelineResult.Namespace, pipelineResult.Name)
	return EtcdSetJson(pipelinePath + "/result", pipelineResult)
}

func GetPipelineResult(pipelineResult *models.PipelineResult) error {
	pipelinePath := fmt.Sprintf(vessel_pipeline_etcd_path, pipelineResult.Namespace, pipelineResult.Name)
	return EtcdGetJson(pipelinePath + "/result",pipelineResult)
}