package pipeline

import (
	"encoding/json"
	"fmt"
	"github.com/containerops/vessel/models"
	"github.com/containerops/vessel/module/etcd"
	"golang.org/x/net/context"
	"log"
	"strconv"
	"strings"
	"time"

	kubeclient "github.com/containerops/vessel/module/kubernetes"
	"github.com/coreos/etcd/client"
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

	// search all stage, start stage which from is ""
	for _, stageVersion := range stageVersions {
		if stageVersion.Name != "" {
			from, err := etcd.GetCurrentStageVersionFromRelation(stageVersion)
			if err != nil {
				// error when get from relation from etcd ,all stage should stop and return err
				stageVersionStateChan <- err.Error()
				notifyBootDone <- true
				continue
			} else if from == "" {
				bootChan <- stageVersion
			}
		}
	}

	runResult := <-stageVersionStateChan
	err = pipelineVersion.Done()
	if err != nil {
		log.Println("error when update pipelineVersion state", err)
	}
	log.Println("all job is done!")
	return runResult, nil
}

// receive bootChan's message start give stage by stage path in etcd
func bootStage(bootChan chan *models.StageVersion, finishChan chan models.StageVersionState, notifyBootDone chan bool) {
	bootMap := make(map[string]bool)
	for {
		select {
		case stageVersion := <-bootChan:
			if _, ok := bootMap[stageVersion.Name]; !ok {
				log.Println("start boot :", stageVersion.Name)
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
			notifyBootDone <- true
			// stageVersionStateChan <- strings.Join(failedList, ",")
			stageVersionStateStr, _ := json.Marshal(stageVersionStateList)
			stageVersionStateChan <- string(stageVersionStateStr)
			return
		}

		stageVersionState := <-finishChan
		stageVersionState.ChangeStageVersionState()
		finishStageNum++
		stageVersionStateList = append(stageVersionStateList, stageVersionState)
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

	count := 0
	sum := len(stageNameMap)

	// pre declare current stageVersion's runState
	var stageVersionState models.StageVersionState
	stageVersionState.PipelineId = stageVersion.PipelineId
	stageVersionState.PipelineVersionId = stageVersion.PipelineVersionId
	stageVersionState.StageId = stageVersion.StageId
	stageVersionState.StageVersionId = stageVersion.PipelineVersionId
	stageVersionState.StageName = stageVersion.Name

	// start watch current stageVersion from stageVersion's state and check thoes stageVersion's running state
	// just in case a stageVersion change after check it's state and before watch it's state
	fromStageVersionStateChan := make(chan string, 255)
	fromStageVersionChan := make(chan *models.StageVersion, 255)
	fromStageVersionStartDoneMap := make(map[string]string)

	for _, fromStageVersionName := range stageNameMap {
		var tempFromStageVersion models.StageVersion
		tempFromStageVersion = *stageVersion
		fromStageVersion := &tempFromStageVersion
		fromStageVersion.Name = fromStageVersionName
		fromStageVersionChan <- fromStageVersion
	}

	watcher := etcd.GetStageVersionFromStageVersionsWatcher(stageVersion)
	go watcheStageVersionState(watcher, fromStageVersionStateChan)

	for {
		if count == sum {
			break
		}
		select {
		case fromStageVersion := <-fromStageVersionChan:
			// get a fromStageVersion from chan to check it's state
			go getCurrentStageVersionState(fromStageVersion, fromStageVersionStateChan)
		case stateInfo := <-fromStageVersionStateChan:
			// get a fromStageVersion state info
			// first check is this state from curren stageVersion's fromStageVersion
			// then check is the state is stage running success or failed
			// 	if state is success,record this and continue wait
			// 	if state is failed,failed current stageVersion
			if len(strings.Split(stateInfo, ",")) == 2 {
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
						stageVersionState.RunResult = StateFailed
						stageVersionState.Detail = "pre stage" + stageName + " is failed"
						changeStageVersionState(stageVersion, bootChan, stageVersionState.RunResult, stageVersionState.Detail)
						finishChan <- stageVersionState
						return
					}
				}
			}
		}
	}

	// start run stage
	stageStartFinish := make(chan models.StageVersionState, 1)
	go startStageInK8S(stageStartFinish, stageVersionState)
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
	startResult := <-stageStartFinish
	stageVersionState = startResult

	changeStageVersionState(stageVersion, bootChan, stageVersionState.RunResult, stageVersionState.Detail)
	finishChan <- stageVersionState
}

func getCurrentStageVersionState(stageVresion *models.StageVersion, stateChan chan string) {
	state, err := etcd.GetCurrentStageVersionState(stageVresion)
	if err != nil {
		log.Println("[getCurrentStageVersionState]:error when get current stageVersion State info :", err)
	} else {
		stateChan <- state
	}
}

func watcheStageVersionState(watcher client.Watcher, fromStageVersionChan chan string) {
	for {
		res, err := watcher.Next(context.Background())
		if err != nil {
			log.Println("[watcheStageVersionState]:error watch stages:", err)
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
		log.Println("[changeStageVersionState]:error when shoutdown stage version:", err)
	}

	if state == StateSuccess || state == StateFailed {
		for _, toStageVersionName := range strings.Split(toStageVersions, ",") {
			if toStageVersionName != "" {
				var toStageVersion models.StageVersion
				toStageVersion = *stageVersion
				toStageVersion.Name = toStageVersionName
				bootChan <- &toStageVersion
			}
		}
	}

}

func startStageInK8S(runResultChan chan models.StageVersionState, runResult models.StageVersionState) {
	pipelineVersion := models.GetPipelineVersion(runResult.PipelineVersionId)
	pipelineSpecTemplate := new(models.PipelineSpecTemplate)
	err := json.Unmarshal([]byte(pipelineVersion.Detail), pipelineSpecTemplate)
	if err != nil {
		log.Printf("Unmarshal PipelineSpecTemplate err : %v\n")
	}
	fmt.Println(pipelineSpecTemplate)
	k8sCh := make(chan string)
	bsCh := make(chan bool)

	go kubeclient.WatchPipelineStatus(pipelineSpecTemplate, kubeclient.Added, k8sCh)

	// runResult.RunResult = <-k8sCh
	if err := kubeclient.StartPipeline(pipelineSpecTemplate); err != nil {
		log.Printf("Start k8s resource pipeline name :%v err : %v\n", pipelineSpecTemplate.MetaData.Name, err)
	}
	go kubeclient.GetPipelineBussinessRes(pipelineSpecTemplate, bsCh)
	fmt.Println("11111111111111")
	k8sRes := ""
	bsRes := true
	for i := 0; i < 2; i++ {
		select {
		case k8sRes = <-k8sCh:
			fmt.Printf("k8sCh return %v\n", k8sRes)
		case bsRes = <-bsCh:
			fmt.Printf("bsCh return %v\n", bsRes)
		}
	}

	if k8sRes == StartFailed {
		fmt.Printf("k8s res %v\n", StateFailed)
		runResult.RunResult = StartFailed
	}
	if k8sRes == StartSucessful {
		fmt.Printf("k8s res %v\n", StartSucessful)
		if bsRes == true {
			fmt.Printf("bs res %v\n", StartSucessful)
			runResult.RunResult = StartSucessful
		} else {
			fmt.Printf("bs res %v\n", StartFailed)
			runResult.RunResult = StartFailed
		}
	}
	fmt.Printf("k8s & bs res %v\n", StartTimeout)
	runResult.RunResult = StartTimeout
	runResult.Detail = StartTimeout

	runResultChan <- runResult
}
