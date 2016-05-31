package etcd

import (
	"fmt"
	"errors"
	"time"

	"github.com/containerops/vessel/models"
	"strings"
)

var (
	VESSEL_PIPELINE_ETCD_PATH = "containerops/vessel/ns_%v/pl_%v"
)

func SavePipeline(info *models.Pipeline) error {
	stagePath := fmt.Sprintf(VESSEL_PIPELINE_ETCD_PATH, info.Namespace, info.Name)
	EtcdSet(stagePath + "/name", info.Name)
	EtcdSet(stagePath + "/namespace", info.Namespace)
	EtcdSet(stagePath + "/stages", strings.Join(info.Stages, ","))
	EtcdSet(stagePath + "/creationTimestamp", info.CreationTimestamp)
	EtcdSet(stagePath + "/deletionTimestamp", info.DeletionTimestamp)
	EtcdSet(stagePath + "/timeoutDuration", info.TimeoutDuration)
	return EtcdSet(stagePath + "/status", info.Status)
}

func GetPipeline(info *models.Pipeline) error {
	stagePath := fmt.Sprintf(VESSEL_PIPELINE_ETCD_PATH, info.Namespace, info.Name)
	response, err := EtcdGet(stagePath)
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
			info.TimeoutDuration = v.Value
		case "/status":
			info.Status = v.Value
		}
	}
	return nil
}

func SetCreationTimestamp(info *models.Pipeline) error {
	stagePath := fmt.Sprintf(VESSEL_PIPELINE_ETCD_PATH, info.Namespace, info.Name)
	info.CreationTimestamp = time.Now().Format("2006-01-02 15:04:05")
	return EtcdSet(stagePath + "/creationTimestamp", info.CreationTimestamp)
}

func SetDeletionTimestamp(info *models.Pipeline) error {
	stagePath := fmt.Sprintf(VESSEL_PIPELINE_ETCD_PATH, info.Namespace, info.Name)
	info.DeletionTimestamp = time.Now().Format("2006-01-02 15:04:05")
	return EtcdSet(stagePath + "/creationTimestamp", info.DeletionTimestamp)
}

func SetPipelineStatusTTL(info *models.Pipeline, timeLife uint64) error {
	stagePath := fmt.Sprintf(VESSEL_PIPELINE_ETCD_PATH, info.Namespace, info.Name)
	return EtcdSetTTL(stagePath + "/status", info.Status, timeLife)
}

func ChangePipelineStatus(info *models.Pipeline) error {
	stagePath := fmt.Sprintf(VESSEL_PIPELINE_ETCD_PATH, info.Namespace, info.Name)
	return EtcdSet(stagePath + "/status", info.Status)
}

func GetPipelineStatus(info *models.Pipeline) (string, error) {
	stagePath := fmt.Sprintf(VESSEL_PIPELINE_ETCD_PATH, info.Namespace, info.Name)
	response, err := EtcdGet(stagePath + "/status")
	if err != nil {
		return "", err
	}
	info.Status = response.Node.Value
	return info.Status, nil
}