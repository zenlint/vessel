package etcd

//coding base models to biz logic

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/containerops/vessel/models"
	"github.com/coreos/etcd/client"
)

const (
	StateNotStart = "not start"
	StateStarting = "working"
	StateSuccess  = "success"
	StateFailed   = "failed"
)

var (
	EtcdClient                        client.Client
	DEFAULT_PIPELINE_ETCD_PATH        = "/containerops/vessel/ws-%d/pj-%d/pl-%d"
	DEFAULT_PIPELINEVERSION_ETCD_PATH = "/containerops/vessel/ws-%d/pj-%d/pl-%d/version/plv-%d"

	START_STAGE_VERSION_CHAN chan bool
)

func init() {
	START_STAGE_VERSION_CHAN = make(chan bool, 1)
}

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
	pipelinePath := fmt.Sprintf(DEFAULT_PIPELINE_ETCD_PATH, pl.WorkspaceId, pl.ProjectId, pl.Id)
	pipelineInfoPath := pipelinePath + "/info"

	EtcdSet(pipelineInfoPath+"/id", strconv.FormatInt(pl.Id, 10))
	EtcdSet(pipelineInfoPath+"/workspaceId", strconv.FormatInt(pl.WorkspaceId, 10))
	EtcdSet(pipelineInfoPath+"/projectId", strconv.FormatInt(pl.ProjectId, 10))
	EtcdSet(pipelineInfoPath+"/name", pl.Name)
	EtcdSet(pipelineInfoPath+"/selfLink", pl.SelfLink)
	EtcdSet(pipelineInfoPath+"/annotations", pl.Annotations)
	EtcdSet(pipelineInfoPath+"/detail", pl.Detail)
}

func SaveStageInfo(stage *models.Stage) {
	// save stage info to etcd
	pipelinePath := fmt.Sprintf(DEFAULT_PIPELINE_ETCD_PATH, stage.WorkspaceId, stage.ProjectId, stage.PipelineId)
	stagePath := pipelinePath + "/stage/" + stage.Name

	EtcdSet(stagePath+"/id", strconv.FormatInt(stage.Id, 10))
	EtcdSet(stagePath+"/workspaceId", strconv.FormatInt(stage.WorkspaceId, 10))
	EtcdSet(stagePath+"/projectId", strconv.FormatInt(stage.ProjectId, 10))
	EtcdSet(stagePath+"/pipelineId", strconv.FormatInt(stage.PipelineId, 10))
	EtcdSet(stagePath+"/name", stage.Name)
	EtcdSet(stagePath+"/detail", stage.Detail)
	EtcdSet(stagePath+"/from", strings.Join(stage.From, ","))
	EtcdSet(stagePath+"/to", strings.Join(stage.To, ","))
}

func SavePipelineVersionInfo(plv *models.PipelineVersion) {
	// save pipeline version info to etcd
	pipelineVersionPath := fmt.Sprintf(DEFAULT_PIPELINEVERSION_ETCD_PATH, plv.WorkspaceId, plv.ProjectId, plv.PipelineId, plv.Id)
	pipelineVersionInfoPath := pipelineVersionPath + "/info"
	pipelineVersionStatePath := pipelineVersionPath + "/state"

	EtcdSet(pipelineVersionInfoPath+"/id", strconv.FormatInt(plv.Id, 10))
	EtcdSet(pipelineVersionInfoPath+"/workspaceId", strconv.FormatInt(plv.WorkspaceId, 10))
	EtcdSet(pipelineVersionInfoPath+"/projectId", strconv.FormatInt(plv.ProjectId, 10))
	EtcdSet(pipelineVersionInfoPath+"/pipelineId", strconv.FormatInt(plv.PipelineId, 10))
	EtcdSet(pipelineVersionInfoPath+"/namespace", plv.Namespace)
	EtcdSet(pipelineVersionInfoPath+"/selfLink", plv.SelfLink)
	EtcdSet(pipelineVersionInfoPath+"/labels", plv.Labels)
	EtcdSet(pipelineVersionInfoPath+"/annotations", plv.Annotations)
	EtcdSet(pipelineVersionInfoPath+"/detail", plv.Detail)
	EtcdSet(pipelineVersionInfoPath+"/stageVersions", plv.StageVersions)

	EtcdSet(pipelineVersionStatePath, plv.Status)
}

func SaveStageVersionInfo(stageVersion *models.StageVersion) {
	// save pipeline version info to etcd
	pipelineVersionPath := fmt.Sprintf(DEFAULT_PIPELINEVERSION_ETCD_PATH, stageVersion.WorkspaceId, stageVersion.ProjectId, stageVersion.PipelineId, stageVersion.PipelineVersionId)
	stageVersionInfoPath := pipelineVersionPath + "/stagev-" + strconv.FormatInt(stageVersion.PipelineVersionId, 10) + "/" + stageVersion.Name + "/info"
	stageVersionStatePath := pipelineVersionPath + "/stagev-" + strconv.FormatInt(stageVersion.PipelineVersionId, 10) + "/" + stageVersion.Name + "/state"

	EtcdSet(stageVersionInfoPath+"/id", strconv.FormatInt(stageVersion.Id, 10))
	EtcdSet(stageVersionInfoPath+"/workspaceId", strconv.FormatInt(stageVersion.WorkspaceId, 10))
	EtcdSet(stageVersionInfoPath+"/projectId", strconv.FormatInt(stageVersion.ProjectId, 10))
	EtcdSet(stageVersionInfoPath+"/pipelineId", strconv.FormatInt(stageVersion.PipelineId, 10))
	EtcdSet(stageVersionInfoPath+"/pipelineVersionId", strconv.FormatInt(stageVersion.PipelineVersionId, 10))
	EtcdSet(stageVersionInfoPath+"/stageId", strconv.FormatInt(stageVersion.StageId, 10))
	EtcdSet(stageVersionInfoPath+"/name", stageVersion.Name)
	EtcdSet(stageVersionInfoPath+"/detail", stageVersion.Detail)

	EtcdSet(stageVersionStatePath+"/workspaceId", strconv.FormatInt(stageVersion.State.WorkspaceId, 10))
	EtcdSet(stageVersionStatePath+"/projectId", strconv.FormatInt(stageVersion.State.ProjectId, 10))
	EtcdSet(stageVersionStatePath+"/pipelineId", strconv.FormatInt(stageVersion.State.PipelineId, 10))
	EtcdSet(stageVersionStatePath+"/pipelineVersionId", strconv.FormatInt(stageVersion.State.PipelineVersionId, 10))
	EtcdSet(stageVersionStatePath+"/stageId", strconv.FormatInt(stageVersion.State.StageId, 10))
	EtcdSet(stageVersionStatePath+"/stageVersionId", strconv.FormatInt(stageVersion.State.StageVersionId, 10))
	EtcdSet(stageVersionStatePath+"/name", stageVersion.State.StageName)
	EtcdSet(stageVersionStatePath+"/runResult", stageVersion.State.RunResult)
	EtcdSet(stageVersionStatePath+"/detail", stageVersion.State.Detail)
}

// get current stageVersion's from relation
func GetCurrentStageVersionFromRelation(stageVersion *models.StageVersion) (string, error) {
	// save stage info to etcd
	pipelinePath := fmt.Sprintf(DEFAULT_PIPELINE_ETCD_PATH, stageVersion.WorkspaceId, stageVersion.ProjectId, stageVersion.PipelineId)
	stagePath := pipelinePath + "/stage/" + stageVersion.Name
	from, err := EtcdGet(stagePath + "/from")
	if err != nil {
		fmt.Println("[IsCurrentStageVersionFromRelationIsNil]:error when get stage from relation :", err.Error())
		return "", err
	}

	return from.Node.Value, nil
}

func StartCurrentStageVersion(stageVersion *models.StageVersion) bool {
	// to make sure this func will not run twice or more at one time
	START_STAGE_VERSION_CHAN <- true
	defer func() {
		<-START_STAGE_VERSION_CHAN
	}()

	pipelineVersionPath := fmt.Sprintf(DEFAULT_PIPELINEVERSION_ETCD_PATH, stageVersion.WorkspaceId, stageVersion.ProjectId, stageVersion.PipelineId, stageVersion.PipelineVersionId)
	stageVersionStatePath := pipelineVersionPath + "/stagev-" + strconv.FormatInt(stageVersion.PipelineVersionId, 10) + "/" + stageVersion.Name + "/state"

	state, _ := EtcdGet(stageVersionStatePath + "/runResult")
	if state != nil {
		runState := strings.Split(state.Node.Value, ",")[1]
		if runState != StateNotStart {
			return false
		}
	}

	EtcdSet(stageVersionStatePath+"/runResult", stageVersion.Name+","+StateStarting)

	return true
}

func GetStageVersionFromStageVersionsWatcher(stageVersion *models.StageVersion) client.Watcher {
	pipelineVersionPath := fmt.Sprintf(DEFAULT_PIPELINEVERSION_ETCD_PATH, stageVersion.WorkspaceId, stageVersion.ProjectId, stageVersion.PipelineId, stageVersion.PipelineVersionId)
	stageVersionPath := pipelineVersionPath + "/stagev-" + strconv.FormatInt(stageVersion.PipelineVersionId, 10) + "/"

	return EtcdWatch(stageVersionPath)
}

func GetCurrentStageVersionState(stageVersion *models.StageVersion) (string, error) {
	pipelineVersionPath := fmt.Sprintf(DEFAULT_PIPELINEVERSION_ETCD_PATH, stageVersion.WorkspaceId, stageVersion.ProjectId, stageVersion.PipelineId, stageVersion.PipelineVersionId)
	stageVersionStatePath := pipelineVersionPath + "/stagev-" + strconv.FormatInt(stageVersion.PipelineVersionId, 10) + "/" + stageVersion.Name + "/state/runResult"
	result, err := EtcdGet(stageVersionStatePath)
	if result != nil {
		return result.Node.Value, nil
	}
	return "", err
}

func ChangeCurrentStageVresionState(stageVersion *models.StageVersion, state, reason string) {
	pipelineVersionPath := fmt.Sprintf(DEFAULT_PIPELINEVERSION_ETCD_PATH, stageVersion.WorkspaceId, stageVersion.ProjectId, stageVersion.PipelineId, stageVersion.PipelineVersionId)
	stageVersionStatePath := pipelineVersionPath + "/stagev-" + strconv.FormatInt(stageVersion.PipelineVersionId, 10) + "/" + stageVersion.Name + "/state"

	EtcdSet(stageVersionStatePath+"/runResult", stageVersion.Name+","+state)
	EtcdSet(stageVersionStatePath+"/detail", reason)
}

func GetCurrentStageVersionToRelation(stageVersion *models.StageVersion) (string, error) {
	// save stage info to etcd
	pipelinePath := fmt.Sprintf(DEFAULT_PIPELINE_ETCD_PATH, stageVersion.WorkspaceId, stageVersion.ProjectId, stageVersion.PipelineId)
	stagePath := pipelinePath + "/stage/" + stageVersion.Name

	to, err := EtcdGet(stagePath + "/to")
	if err != nil {
		fmt.Println("[IsCurrentStageVersiontoRelationIsNil]:error when get stage to relation :", err.Error())
		return "", err
	}

	return to.Node.Value, nil
}

// ///////////////////////////////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////////////////////////////

// func SavePipelineId(path, id string) {
// 	EtcdSet(path, id)
// }

// func GetStageNamesByPipelineVersion(pipelineVersion *models.PipelineVersion) []string {
// 	pipelinePath := fmt.Sprintf(DEFAULT_PIPELINE_ETCD_PATH, pipelineVersion.WorkspaceId, pipelineVersion.ProjectId, pipelineVersion.PipelineId)
// 	// pipelineVersionPath := fmt.Sprintf(DEFAULT_PIPELINEVERSION_ETCD_PATH, pipelineVersion.WorkspaceId, pipelineVersion.ProjectId, pipelineVersion.Id, pipelineVersion.Id)
// 	stageList, _ := EtcdGet(pipelinePath + "/allstage")
// 	if stageList != nil {
// 		return strings.Split(stageList.Node.Value, ",")
// 	}
// 	return make([]string, 0)
// }

// func GetStageFromInfoByPipelineVersionAndStageName(pipelineVersion *models.PipelineVersion, stageName string) (string, string) {
// 	pipelinePath := fmt.Sprintf(DEFAULT_PIPELINE_ETCD_PATH, pipelineVersion.WorkspaceId, pipelineVersion.ProjectId, pipelineVersion.PipelineId)
// 	pipelineVersionPath := fmt.Sprintf(DEFAULT_PIPELINEVERSION_ETCD_PATH, pipelineVersion.WorkspaceId, pipelineVersion.ProjectId, pipelineVersion.Id, pipelineVersion.Id)
// 	stagePath := pipelinePath + stageName
// 	stageVersionPath := pipelineVersionPath + stageName
// 	fromInfo, _ := EtcdGet(stagePath + "/from")
// 	if fromInfo != nil {
// 		return stageVersionPath, fromInfo.Node.Value
// 	}
// 	return stageVersionPath, ""
// }

// func GetStageVersionInfoByPath(stageVersionStagePath string) (string, *models.StageVersion) {
// 	stageVersion := new(models.StageVersion)

// 	stageName := stageVersionStagePath[strings.LastIndex(stageVersionStagePath, "/")+1:]
// 	// stageVersionPath
// 	stageVersionPath := stageVersionStagePath[:strings.LastIndex(stageVersionStagePath, "/")]
// 	// pipelineVersionPath
// 	pipelineVersionPath := stageVersionPath[:strings.LastIndex(stageVersionPath, "/")]
// 	// pipelineVersionId
// 	pipelineVersionID := pipelineVersionPath[strings.LastIndex(pipelineVersionPath, "-")+1:]
// 	// pipelineID
// 	pipelineIDInfo, _ := EtcdGet(pipelineVersionPath + "/pipelineId")
// 	pipelineID := pipelineIDInfo.Node.Value
// 	// stagePath
// 	stagePath := pipelineVersionPath[:strings.LastIndex(pipelineVersionPath, "/")] + "/pl-" + pipelineID + "/stage"
// 	// stageID
// 	stageIDInfo, _ := EtcdGet(stagePath + "/" + stageName + "/id")
// 	stageID := stageIDInfo.Node.Value
// 	// get current stage from info
// 	fromStageNamesInfo, _ := EtcdGet(stagePath + "/" + stageName + "/from")
// 	fromStageNames := fromStageNamesInfo.Node.Value

// 	// get current stage to info
// 	toStageNamesInfo, _ := EtcdGet(stagePath + "/" + stageName + "/to")
// 	toStageNames := toStageNamesInfo.Node.Value

// 	stageVersion.PipelineId, _ = strconv.ParseInt(pipelineID, 10, 64)
// 	// strconv.ParseInt(s string, base int, bitSize int)
// 	stageVersion.PipelineVersionId, _ = strconv.ParseInt(pipelineVersionID, 10, 64)
// 	stageVersion.StageId, _ = strconv.ParseInt(stageID, 10, 64)
// 	stageVersion.Name = stageName
// 	stageVersion.From = strings.Split(fromStageNames, ",")
// 	stageVersion.To = strings.Split(toStageNames, ",")

// 	stateInfo, _ := EtcdGet(stageVersionPath + "/" + stageVersion.Name + "/state")
// 	state := ""
// 	if stateInfo != nil {
// 		state = stateInfo.Node.Value
// 	}

// 	return state, stageVersion
// }
