package pipeline

import (
	"encoding/json"
	"github.com/containerops/vessel/models"
	"github.com/containerops/vessel/module/etcd"
	"golang.org/x/net/context"
	log "github.com/golang/glog"
	"strconv"
	"strings"
	"time"

	kubeclient "github.com/containerops/vessel/module/kubernetes"
	"github.com/coreos/etcd/client"
	"errors"
)

const (
	StateNotStart  = "not start"
	StateStarting  = "working"
	StateSuccess   = "success"
	StateFailed    = "failed"
	StartSucessful = "OK"
	StartFailed    = "ERROR"
	StartTimeout   = "TIMEOUT"
)

// BootPipelineVersion : start a pipelineVersion,boot the stage and return the result
func BootPipelineVersion(pipelineId int64) (string, error) {
	// get pipeline info by pipelineId
	pipelineInfo, err := new(models.Pipeline).GetPipelineInfoByPipelineId(pipelineId)
	if err != nil {
		return "", err
	}
	// storage pipeline info to etcd
	etcd.SavePipelineInfo(pipelineInfo)

	// get stage infos by pipeline info
	stages, err := new(models.Stage).GetStagesByPipelineInfo(pipelineInfo)
	if err != nil {
		return "", err
	}
	// storage stage infos to etcd
	for _, stage := range stages {
		etcd.SaveStageInfo(stage)
	}

	// generate new pipelineVersion and save it to db and set state to running
	pipelineVersion := new(models.PipelineVersion)
	pipelineVersion.Id = time.Now().UnixNano()
	pipelineVersion.WorkspaceId = pipelineInfo.WorkspaceId
	pipelineVersion.ProjectId = pipelineInfo.ProjectId
	pipelineVersion.PipelineId = pipelineInfo.Id
	pipelineVersion.Namespace = "plv" + "-" + strconv.FormatInt(pipelineVersion.Id, 10)
	pipelineVersion.SelfLink = ""
	pipelineVersion.Labels = pipelineInfo.Labels
	pipelineVersion.Annotations = pipelineInfo.Annotations
	pipelineVersion.Detail = pipelineInfo.Detail
	pipelineVersion.StageVersions = strconv.FormatInt(pipelineVersion.Id, 10)
	pipelineVersion.Status = StateStarting

	err = pipelineVersion.Save()
	if err != nil {
		return "", err
	}
	// save pipelineVersion info to etcd
	etcd.SavePipelineVersionInfo(pipelineVersion)

	// generate stageVersion info and save it to db and etcd
	// and record to a slice
	stageVersions := make([]*models.StageVersion, 0)
	for _, stage := range stages {
		stageVersion := new(models.StageVersion)
		stageVersion.Id = time.Now().UnixNano()
		stageVersion.WorkspaceId = stage.WorkspaceId
		stageVersion.ProjectId = stage.ProjectId
		stageVersion.PipelineId = stage.PipelineId
		stageVersion.PipelineVersionId = pipelineVersion.Id
		stageVersion.StageId = stage.Id
		stageVersion.Name = stage.Name
		stageVersion.Detail = stage.Detail

		runState := new(models.StageVersionState)
		runState.WorkspaceId = stageVersion.WorkspaceId
		runState.ProjectId = stageVersion.ProjectId
		runState.PipelineId = stageVersion.PipelineId
		runState.PipelineVersionId = stageVersion.PipelineVersionId
		runState.StageId = stageVersion.StageId
		runState.StageVersionId = stageVersion.Id
		runState.StageName = stageVersion.Name
		runState.RunResult = stageVersion.Name + "," + StateNotStart
		runState.Detail = stageVersion.Detail

		stageVersion.State = runState

		err = stageVersion.Save()
		if err != nil {
			return "", err
		}

		etcd.SaveStageVersionInfo(stageVersion)

		stageVersions = append(stageVersions, stageVersion)
	}

	// start to boot stage
	bootChan := make(chan *models.StageVersion, 100)
	finishChan := make(chan models.StageVersionState, 100)
	stageVersionStateChan := make(chan string, 1)
	notifyBootDone := make(chan bool, 1)

	go bootStage(bootChan, finishChan, notifyBootDone)
	go isFinish(finishChan, stageVersionStateChan, len(stageVersions), notifyBootDone)

	// search all stage, start stage which from is empty
	for _, stageVersion := range stageVersions {
		if stageVersion.Name != "" {
			from, err := etcd.GetCurrentStageVersionFromRelation(stageVersion)
			if err != nil {
				// error when get from relation from etcd ,all stage should stop and return err
				stageVersionStateChan <- err.Error()
				notifyBootDone <- true
				break
			} else if from == "" {
				bootChan <- stageVersion
			}
		}
	}

	stageVersionState := <-stageVersionStateChan
	err = pipelineVersion.Done()
	if err != nil {
		log.Infoln("error when update pipelineVersion state", err)
	}
	log.Infoln("all job is done!")
	return stageVersionState, nil
}

// receive bootChan's message start give stage by stage path in etcd
func bootStage(bootChan chan *models.StageVersion, finishChan chan models.StageVersionState, notifyBootDone chan bool) {
	bootMap := make(map[string]bool)
	for {
		select {
		case stageVersion := <-bootChan:
			if _, ok := bootMap[stageVersion.Name]; !ok {
				log.Infoln("start boot :", stageVersion.Name)
				bootMap[stageVersion.Name] = true
				go startStage(stageVersion, bootChan, finishChan)
			}
		case <-notifyBootDone:
			return
		}
	}
}

// count finish stage num,send single when all stage is start
func isFinish(finishChan chan models.StageVersionState, stageVersionStateChan chan string, sumStage int, notifyBootDone chan bool) {
	finishStageNum := 0
	stageVersionStateList := make([]models.StageVersionState, 0)

	for {
		if finishStageNum == sumStage {
			log.Infoln("finishStageNum == sumStage")
			notifyBootDone <- true

			// stageVersionStateChan <- strings.Join(failedList, ",")
			stageVersionStateStr, _ := json.Marshal(stageVersionStateList)
			stageVersionStateChan <- string(stageVersionStateStr)
			return
		}
		log.Infoln("finishStageNum = ",finishStageNum)
		stageVersionState := <-finishChan
		stageVersionState.ChangeStageVersionState()
		stageVersionStateList = append(stageVersionStateList, stageVersionState)
		if stageVersionState.RunResult != StartSucessful{
			finishStageNum = sumStage
		}else{
			finishStageNum++
		}
		log.Infoln("finishStageNum++ = ",finishStageNum)
	}
}

// stage stage by given name,after start give single to finishChan
func startStage(stageVersion *models.StageVersion, bootChan chan *models.StageVersion, finishChan chan models.StageVersionState) {
	// try to start current stageVersion if start success ,return true,else return false
	if !etcd.StartCurrentStageVersion(stageVersion) {
		return
	}

	// get all stage name that current stageVersion rely
	stageNameMap := make(map[string]string)
	fromStageRelation, _ := etcd.GetCurrentStageVersionFromRelation(stageVersion)
	for _, stageName := range strings.Split(fromStageRelation, ",") {
		if stageName != "" {
			stageNameMap[stageName] = stageName
		}
	}

	// start watch current stageVersion from stageVersion's state and check thoes stageVersion's running state
	// just in case a stageVersion change after check it's state and before watch it's state
	fromStageVersionStateChan := make(chan string, 255)
	fromStageVersionStartDoneMap := make(map[string]string)

	for _, fromStageVersionName := range stageNameMap {
		fromStageVersion := *stageVersion
		fromStageVersion.Name = fromStageVersionName
		go getCurrentStageVersionState(&fromStageVersion, fromStageVersionStateChan)
	}

	watcher := etcd.GetStageVersionFromStageVersionsWatcher(stageVersion)
	go watchStageVersionState(watcher, fromStageVersionStateChan)

	for count, sum := 0,len(stageNameMap); count < sum; {
		select {
		case stateInfo := <-fromStageVersionStateChan:
			// get a fromStageVersion state info
			// first check is this state from curren stageVersion's fromStageVersion
			// then check is the state is stage running success or failed
			// 	if state is success,record this and continue wait
			// 	if state is failed,failed current stageVersion
			infoArr := strings.Split(stateInfo, ",")
			if len(infoArr) != 2 {
				continue;
			}
			stageName := strings.Split(stateInfo, ",")[0]
			stageState := strings.Split(stateInfo, ",")[1]
			if _, exist := stageNameMap[stageName]; exist {
				if stageState == StateSuccess {
					fromStageVersionStartDoneMap[stageName] = stageName
					count = len(fromStageVersionStartDoneMap)
				} else if stageState == StateFailed {
					// if got a failed stageVersion,shutdown current stageVersion
					// notify current stageVresion's to stageVersion
					// return func

					stageVersionState := formatStageVersionState(stageVersion, errors.New("pre stage" + stageName + " is failed"))
					changeStageVersionState(stageVersion, bootChan, stageVersionState.RunResult, stageVersionState.Detail)
					finishChan <- *stageVersionState
					return
				}
			}
		}
	}

	// start run stage
	err := startStageInK8S(stageVersion.PipelineVersionId,stageVersion.Name)
	stageVersionState := formatStageVersionState(stageVersion, err)

	/*timeout := make(chan bool, 1)
	// stageStartFinish <- stageVersionState
	go startStageInK8S(stageStartFinish, stageVersionState)
	go func() {
		d, _ := time.ParseDuration("5s")
		time.Sleep(d)
		timeout <- true
	}()

	// check is stage run timeout,and get run result
	select {
	case <-timeout:
		stageVersionState.RunResult = StateFailed
		stageVersionState.Detail = "time out"
	case startResult := <-stageStartFinish:
		stageVersionState = startResult
	}
	*/
	changeStageVersionState(stageVersion, bootChan, stageVersionState.RunResult, stageVersionState.Detail)
	finishChan <- *stageVersionState
}

func getCurrentStageVersionState(stageVersion *models.StageVersion, stateChan chan string) {
	state, err := etcd.GetCurrentStageVersionState(stageVersion)
	if err != nil {
		log.Infoln("[getCurrentStageVersionState]:error when get current stageVersion State info :", err)
	} else {
		stateChan <- state
	}
}

func watchStageVersionState(watcher client.Watcher, fromStageVersionChan chan string) {
	for {
		res, err := watcher.Next(context.Background())
		if err != nil {
			log.Infoln("[watcheStageVersionState]:error watch stages:", err)
		}

		if res.Action == "set" || res.Action == "update" {
			changeStageInfo := res.Node.Value
			fromStageVersionChan <- changeStageInfo
		}
	}
}

func changeStageVersionState(stageVersion *models.StageVersion, bootChan chan *models.StageVersion, state, reason string) {
	etcd.ChangeCurrentStageVresionState(stageVersion, state, reason)
	toStageVersions, err := etcd.GetCurrentStageVersionToRelation(stageVersion)
	if err != nil {
		log.Infoln("[changeStageVersionState]:error when shoutdown stage version:", err)
	}
	log.Infoln(state)

	if state == StateSuccess || state == StateFailed {
		for _, toStageVersionName := range strings.Split(toStageVersions, ",") {
			log.Infoln(toStageVersions)
			if toStageVersionName != "" {
				log.Infoln(toStageVersionName)
				toStageVersion := *stageVersion
				toStageVersion.Name = toStageVersionName
				bootChan <- &toStageVersion
			}
		}
	}
}

func startStageInK8S(pipelineVersionId int64,stageName string) (err error){
	log.Infoln("Enter startStageInK8S to start stage ", stageName)
	pipelineVersion := models.GetPipelineVersion(pipelineVersionId)
	pipelineSpecTemplate := new(models.PipelineSpecTemplate)

	if err = json.Unmarshal([]byte(pipelineVersion.Detail), pipelineSpecTemplate); err != nil {
		log.Infoln("Unmarshal PipelineSpecTemplate err: ",err)
		return err
	}

	log.Infoln("goting to deal with pipelinePecTemplate detail = ", pipelineSpecTemplate)
	k8sCh := make(chan string)
	bsCh := make(chan bool)

	go kubeclient.WatchPipelineStatus(pipelineSpecTemplate, stageName, kubeclient.Added, k8sCh)

	// runResult.RunResult = <-k8sCh
	err = kubeclient.StartPipeline(pipelineSpecTemplate, stageName)
	if err != nil {
		log.Infoln("Start k8s resource pipeline name: ", pipelineSpecTemplate.MetaData.Name," err: ", err)
	}
	go kubeclient.GetPipelineBussinessRes(pipelineSpecTemplate, bsCh)
	for i := 0; i < 2; i++ {
		select {
		case k8sRes := <-k8sCh:
			log.Infoln("k8sCh start stage name = ", stageName," return ", k8sRes)
		case bsRes := <-bsCh:
			if !bsRes{
				err = errors.New("Get pipeline bussiness Res wrong")
			}
		}
	}
	return err
}

func formatStageVersionState(stageVersion *models.StageVersion,err error) *models.StageVersionState {
	// pre declare current stageVersion's runState
	stageVersionState := models.StageVersionState{
		PipelineId:stageVersion.PipelineId,
		PipelineVersionId: stageVersion.PipelineVersionId,
		StageId: stageVersion.StageId,
		StageVersionId: stageVersion.PipelineVersionId,
		StageName: stageVersion.Name}
	log.Infoln("k8s module stage name = ", stageVersion.Name," ret ", err)
	if err == nil {
		stageVersionState.RunResult = StateSuccess
		stageVersionState.Detail = StateSuccess
	}else if err.Error() == StartTimeout{
		stageVersionState.RunResult = StartTimeout
		stageVersionState.Detail = StartTimeout
	}else{
		stageVersionState.RunResult = StartFailed
		stageVersionState.Detail = err.Error()
	}
	return &stageVersionState
}