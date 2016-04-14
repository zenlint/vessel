package etcd

//coding base models to biz logic

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/containerops/vessel/models"
	"github.com/coreos/etcd/client"
)

var (
	EtcdClient                        client.Client
	DEFAULT_PIPELINE_ETCD_PATH        = "/containerops/vessel/ws-%d/pj-%d/pl-%d/stage/"
	DEFAULT_PIPELINEVERSION_ETCD_PATH = "/containerops/vessel/ws-%d/pj-%d/plv-%d/stagev-%d/"
)

//Sync ETCD
func SyncETCD() error {
	return nil
}

func EtcdSet(key, value string) error {
	return models.EtcdSet(key, value)
}

func EtcdGet(key string) (*client.Response, error) {
	return models.EtcdGet(key)
}

func EtcdWatch(path string) client.Watcher {
	return models.EtcdWatch(path)
}

func SavePipelineInfo(pl *models.Pipeline) {
	// save pipeline info to etcd
	stageNames := make([]string, 0)
	for _, stage := range pl.Stages {
		stageNames = append(stageNames, stage.Name)
	}
	pipelinePath := fmt.Sprintf(DEFAULT_PIPELINE_ETCD_PATH, pl.WorkspaceId, pl.ProjectId, pl.Id)
	EtcdSet(pipelinePath+"/allstage", strings.Join(stageNames, ","))
}

func SaveStageInfo(stage *models.Stage, stagePath string) {
	// save stage info to etcd
	EtcdSet(stagePath+"/id", strconv.FormatInt(stage.Id, 10))
	EtcdSet(stagePath+"/name", stage.Name)
	EtcdSet(stagePath+"/detail", stage.Detail)
	EtcdSet(stagePath+"/from", strings.Join(stage.From, ","))
	EtcdSet(stagePath+"/to", strings.Join(stage.To, ","))
}

func SavePipelineVersionInfo(plv *models.PipelineVersion) {
	// save pipeline version info to etcd
}

func SavePipelineId(path, id string) {
	EtcdSet(path, id)
}

func GetStageNamesByPipelineVersion(pipelineVersion *models.PipelineVersion) []string {
	pipelinePath := fmt.Sprintf(DEFAULT_PIPELINE_ETCD_PATH, pipelineVersion.WorkspaceId, pipelineVersion.ProjectId, pipelineVersion.PipelineId)
	// pipelineVersionPath := fmt.Sprintf(DEFAULT_PIPELINEVERSION_ETCD_PATH, pipelineVersion.WorkspaceId, pipelineVersion.ProjectId, pipelineVersion.Id, pipelineVersion.Id)
	stageList, _ := EtcdGet(pipelinePath + "/allstage")
	if stageList != nil {
		return strings.Split(stageList.Node.Value, ",")
	}
	return make([]string, 0)
}

func GetStageFromInfoByPipelineVersionAndStageName(pipelineVersion *models.PipelineVersion, stageName string) (string, string) {
	pipelinePath := fmt.Sprintf(DEFAULT_PIPELINE_ETCD_PATH, pipelineVersion.WorkspaceId, pipelineVersion.ProjectId, pipelineVersion.PipelineId)
	pipelineVersionPath := fmt.Sprintf(DEFAULT_PIPELINEVERSION_ETCD_PATH, pipelineVersion.WorkspaceId, pipelineVersion.ProjectId, pipelineVersion.Id, pipelineVersion.Id)
	stagePath := pipelinePath + stageName
	stageVersionPath := pipelineVersionPath + stageName
	fromInfo, _ := EtcdGet(stagePath + "/from")
	if fromInfo != nil {
		return stageVersionPath, fromInfo.Node.Value
	}
	return stageVersionPath, ""
}

func GetStageVersionInfoByPath(stageVersionStagePath string) (string, *models.StageVersion) {
	stageVersion := new(models.StageVersion)

	stageName := stageVersionStagePath[strings.LastIndex(stageVersionStagePath, "/")+1:]
	// stageVersionPath
	stageVersionPath := stageVersionStagePath[:strings.LastIndex(stageVersionStagePath, "/")]
	// pipelineVersionPath
	pipelineVersionPath := stageVersionPath[:strings.LastIndex(stageVersionPath, "/")]
	// pipelineVersionId
	pipelineVersionID := pipelineVersionPath[strings.LastIndex(pipelineVersionPath, "-")+1:]
	// pipelineID
	pipelineIDInfo, _ := EtcdGet(pipelineVersionPath + "/pipelineId")
	pipelineID := pipelineIDInfo.Node.Value
	// stagePath
	stagePath := pipelineVersionPath[:strings.LastIndex(pipelineVersionPath, "/")] + "/pl-" + pipelineID + "/stage"
	// stageID
	stageIDInfo, _ := EtcdGet(stagePath + "/" + stageName + "/id")
	stageID := stageIDInfo.Node.Value
	// get current stage from info
	fromStageNamesInfo, _ := EtcdGet(stagePath + "/" + stageName + "/from")
	fromStageNames := fromStageNamesInfo.Node.Value

	// get current stage to info
	toStageNamesInfo, _ := EtcdGet(stagePath + "/" + stageName + "/to")
	toStageNames := toStageNamesInfo.Node.Value

	stageVersion.PipelineId, _ = strconv.ParseInt(pipelineID, 10, 64)
	// strconv.ParseInt(s string, base int, bitSize int)
	stageVersion.PipelineVersionId, _ = strconv.ParseInt(pipelineVersionID, 10, 64)
	stageVersion.StageId, _ = strconv.ParseInt(stageID, 10, 64)
	stageVersion.Name = stageName
	stageVersion.From = strings.Split(fromStageNames, ",")
	stageVersion.To = strings.Split(toStageNames, ",")

	stateInfo, _ := EtcdGet(stageVersionPath + "/" + stageVersion.Name + "/state")
	state := ""
	if stateInfo != nil {
		state = stateInfo.Node.Value
	}

	return state, stageVersion
}
