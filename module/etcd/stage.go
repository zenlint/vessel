package etcd

import (
	"fmt"
	"errors"

	"github.com/containerops/vessel/models"
	"strings"
)

var (
	VESSEL_STAGE_ETCD_PATH = "containerops/vessel/ns_%v/sn_%v"
)

func SaveStage(info *models.Stage) error {
	stagePath := fmt.Sprintf(VESSEL_STAGE_ETCD_PATH, info.Namespace, info.Name)
	EtcdSet(stagePath + "/name", info.Name)
	EtcdSet(stagePath + "/namespace", info.Namespace)
	EtcdSet(stagePath + "/replicas", info.Replicas)
	EtcdSet(stagePath + "/image", info.Image)
	EtcdSet(stagePath + "/envName", info.EnvName)
	EtcdSet(stagePath + "/envValue", info.EnvValue)
	EtcdSet(stagePath + "/port", info.Port)
	EtcdSet(stagePath + "/dependence", strings.Join(info.Dependence,","))
	return EtcdSet(stagePath + "/status", info.Status)
}

func GetStage(info *models.Stage) error {
	stagePath := fmt.Sprintf(VESSEL_STAGE_ETCD_PATH, info.Namespace, info.Name)
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
		case "/replicas":
			info.Replicas = v.Value
		case "/image":
			info.Image = v.Value
		case "/envName":
			info.EnvName = v.Value
		case "/envValue":
			info.EnvValue = v.Value
		case "/port":
			info.Port = v.Value
		case "/dependence":
			info.Dependence = strings.Split(v.Value,",")
		case "/status":
			info.Status = v.Value
		}
	}
	return nil
}

func SetStageStatusTTL(info *models.Stage, timeLife uint64) error {
	stagePath := fmt.Sprintf(VESSEL_STAGE_ETCD_PATH, info.Namespace, info.Name)
	return EtcdSetTTL(stagePath + "/status", info.Status, timeLife)
}

func ChangeStageStatus(info *models.Stage) error {
	stagePath := fmt.Sprintf(VESSEL_STAGE_ETCD_PATH, info.Namespace, info.Name)
	return EtcdSet(stagePath + "/status", info.Status)
}

func GetStageStatus(info *models.Stage) (string, error) {
	stagePath := fmt.Sprintf(VESSEL_STAGE_ETCD_PATH, info.Namespace, info.Name)
	response, err := EtcdGet(stagePath + "/status")
	if err != nil {
		return "", err
	}
	info.Status = response.Node.Value
	return response.Node.Value, nil
}