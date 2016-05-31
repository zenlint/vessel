package etcd

import (
	"fmt"
	"errors"
	"time"
	"strings"
	"strconv"
	"log"

	"github.com/containerops/vessel/models"
)

var (
	vessel_pipeline_etcd_path = "containerops/vessel/ns_%v/pl_%v"
)

func SavePipeline(info *models.Pipeline) error {
	pipelinePath := fmt.Sprintf(vessel_pipeline_etcd_path, info.Namespace, info.Name)
	EtcdSet(pipelinePath + "/name", info.Name)
	EtcdSet(pipelinePath + "/namespace", info.Namespace)
	EtcdSet(pipelinePath + "/stages", strings.Join(info.Stages, ","))
	EtcdSet(pipelinePath + "/creationTimestamp", info.CreationTimestamp)
	EtcdSet(pipelinePath + "/deletionTimestamp", info.DeletionTimestamp)
	EtcdSet(pipelinePath + "/timeoutDuration", strconv.FormatInt(info.TimeoutDuration, 10))
	return EtcdSet(pipelinePath + "/status", info.Status)
}

func GetPipeline(info *models.Pipeline) error {
	pipelinePath := fmt.Sprintf(vessel_pipeline_etcd_path, info.Namespace, info.Name)
	response, err := EtcdGet(pipelinePath)
	if err != nil {
		return err
	}
	if !response.Node.Dir {
		return errors.New(fmt.Sprintf("Stage name : %v was not found in ETCD where namespace = %v", info.Name, info.Namespace))
	}
	for _, v := range response.Node.Nodes {
		switch v.Key {
		case "/name":
			info.Name = v.Value
		case "/namespace":
			info.Namespace = v.Value
		case "/stages":
			info.Stages = strings.Split(v.Value, ",")
		case "/creationTimestamp":
			info.CreationTimestamp = v.Value
		case "/deletionTimestamp":
			info.DeletionTimestamp = v.Value
		case "/timeoutDuration":
			value, err := strconv.ParseInt(v.Value, 10, 0)
			if err != nil {
				log.Println(err)
			} else {
				info.TimeoutDuration = value
			}
		case "/status":
			info.Status = v.Value
		}
	}
	return nil
}

func SetCreationTimestamp(info *models.Pipeline) error {
	pipelinePath := fmt.Sprintf(vessel_pipeline_etcd_path, info.Namespace, info.Name)
	info.CreationTimestamp = time.Now().Format("2006-01-02 15:04:05")
	return EtcdSet(pipelinePath + "/creationTimestamp", info.CreationTimestamp)
}

func SetDeletionTimestamp(info *models.Pipeline) error {
	pipelinePath := fmt.Sprintf(vessel_pipeline_etcd_path, info.Namespace, info.Name)
	info.DeletionTimestamp = time.Now().Format("2006-01-02 15:04:05")
	return EtcdSet(pipelinePath + "/creationTimestamp", info.DeletionTimestamp)
}

func SetPipelineTTL(info *models.Pipeline, timeLife uint64) error {
	pipelinePath := fmt.Sprintf(vessel_pipeline_etcd_path, info.Namespace, info.Name)
	return EtcdSetDirTTL(pipelinePath, timeLife)
}

func ChangePipelineStatus(info *models.Pipeline) error {
	pipelinePath := fmt.Sprintf(vessel_pipeline_etcd_path, info.Namespace, info.Name)
	return EtcdSet(pipelinePath + "/status", info.Status)
}

func GetPipelineStatus(info *models.Pipeline) (string, error) {
	pipelinePath := fmt.Sprintf(vessel_pipeline_etcd_path, info.Namespace, info.Name)
	response, err := EtcdGet(pipelinePath + "/status")
	if err != nil {
		return "", err
	}
	info.Status = response.Node.Value
	return info.Status, nil
}