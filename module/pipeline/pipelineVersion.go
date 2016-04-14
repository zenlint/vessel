package pipeline

import (
	"encoding/json"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/containerops/vessel/models"
	"github.com/containerops/vessel/module/etcd"

	"golang.org/x/net/context"
)

const (
	StageStateNotStart = "not start"
	StageStateStarting = "working"
	StageStateSuccess  = "success"
	StageStateFailed   = "failed"
)

// BootPipelineVersion : start a pipelineVersion,boot the stage and return the result
func BootPipelineVersion(pipelineVersion *models.PipelineVersion) string {
	bootChan := make(chan string, 100)
	finishChan := make(chan models.StageVersionState, 100)
	stageVersionStateChan := make(chan string, 1)
	notifyBootDone := make(chan bool, 1)

	stageNames := etcd.GetStageNamesByPipelineVersion(pipelineVersion)
	sumStage := len(stageNames)

	go bootStage(bootChan, finishChan, notifyBootDone)
	go isFinish(finishChan, stageVersionStateChan, sumStage, notifyBootDone)

	// search all stage, start stage which from is ""
	for _, stageName := range stageNames {
		if stageName != "" {
			stageVersionPath, from := etcd.GetStageFromInfoByPipelineVersionAndStageName(pipelineVersion, stageName)
			if from == "" {
				bootChan <- stageVersionPath
			}
		}
	}

	runResult := <-stageVersionStateChan
	log.Println("all job is done!")
	return runResult
}

// receive bootChan's message start give stage by stage path in etcd
func bootStage(bootChan chan string, finishChan chan models.StageVersionState, notifyBootDone chan bool) {
	bootMap := make(map[string]bool)
	for {
		select {
		case stagePath := <-bootChan:
			if _, ok := bootMap[stagePath]; !ok {
				bootMap[stagePath] = true
				go startStage(stagePath, bootChan, finishChan)
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
			notifyBootDone <- false
			// stageVersionStateChan <- strings.Join(failedList, ",")
			stageVersionStateStr, _ := json.Marshal(stageVersionStateList)
			stageVersionStateChan <- string(stageVersionStateStr)
			return
		}

		stageVersionState := <-finishChan
		finishStageNum++
		stageVersionStateList = append(stageVersionStateList, stageVersionState)
	}
}

// stage stage by given name,after start give single to finishChan
func startStage(stageVersionStagePath string, bootChan chan string, finishChan chan models.StageVersionState) {
	stateInfo, stageVersionInfo := etcd.GetStageVersionInfoByPath(stageVersionStagePath)
	if stateInfo != "" {
		return
	}

	etcd.EtcdSet(stageVersionStagePath[:strings.LastIndex(stageVersionStagePath, "/")]+"/state", stageVersionInfo.Name+","+StageStateNotStart)

	stageNameMap := make(map[string]string)
	for _, stageName := range stageVersionInfo.From {
		if stageName != "" {
			stageNameMap[stageName] = stageName
		}
	}

	count := 0
	sum := len(stageNameMap)
	var stageVersionState models.StageVersionState
	stageVersionState.PipelineId = strconv.FormatInt(stageVersionInfo.PipelineId, 10)
	stageVersionState.PipelineVersionId = strconv.FormatInt(stageVersionInfo.PipelineVersionId, 10)
	stageVersionState.StageId = strconv.FormatInt(stageVersionInfo.StageId, 10)
	stageVersionState.StageVersionId = strconv.FormatInt(stageVersionInfo.PipelineVersionId, 10)
	stageVersionState.StageName = stageVersionInfo.Name

	// check is there is some stage finish start before this stage start
	for _, fromStageName := range stageNameMap {
		fromStageVersionStatePath := stageVersionStagePath[:strings.LastIndex(stageVersionStagePath, "/")] + "/" + fromStageName + "/state"
		fromStageVersionStateInfo, _ := etcd.EtcdGet(fromStageVersionStatePath)
		if fromStageVersionStateInfo != nil {
			info := fromStageVersionStateInfo.Node.Value
			if strings.Split(info, ",")[1] == StageStateSuccess {
				count++
			} else if strings.Split(info, ",")[1] == StageStateFailed {
				for _, v := range stageVersionInfo.To {
					if v != "" {
						bootChan <- stageVersionStagePath[:strings.LastIndex(stageVersionStagePath, "/")] + "/" + v
					}
				}
				stageVersionState.RunResult = StageStateFailed
				stageVersionState.Detail = "pre stage " + strings.Split(info, ",")[0] + " is failed"
				finishChan <- stageVersionState
				etcd.EtcdSet(stageVersionStagePath+"/state", stageVersionInfo.Name+","+StageStateFailed)
				return
			}
		}
	}

	watcher := etcd.EtcdWatch(stageVersionStagePath[:strings.LastIndex(stageVersionStagePath, "/")] + "/")
	for {
		// if all stage dependence is start,exit loop and start current stage
		if count == sum {
			break
		}

		res, err := watcher.Next(context.Background())
		if err != nil {
			log.Println("error watch stages:", err)
		}
		if res.Action == "set" || res.Action == "update" {
			changeStageInfo := res.Node.Value
			if _, ok := stageNameMap[strings.Split(changeStageInfo, ",")[0]]; ok {
				if strings.Split(changeStageInfo, ",")[1] == StageStateSuccess {
					count++
				} else if strings.Split(changeStageInfo, ",")[1] == StageStateFailed {
					for _, v := range stageVersionInfo.To {
						if v != "" {
							bootChan <- stageVersionStagePath[:strings.LastIndex(stageVersionStagePath, "/")] + "/" + v
						}
						// }
					}
					stageVersionState.RunResult = StageStateFailed
					stageVersionState.Detail = "pre stage " + strings.Split(changeStageInfo, ",")[0] + " is failed"
					finishChan <- stageVersionState
					etcd.EtcdSet(stageVersionStagePath+"/state", stageVersionInfo.Name+","+StageStateFailed)
					return
				}
			}
		}
	}

	// start run stage
	stageStartFinish := make(chan models.StageVersionState, 1)
	timeout := make(chan bool, 1)
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
		stageVersionState.RunResult = StageStateFailed
		stageVersionState.Detail = "time out"
	case startResult := <-stageStartFinish:
		stageVersionState = startResult
	}

	// check stage run result and set it to etcd
	if stageVersionState.RunResult == StageStateSuccess {
		etcd.EtcdSet(stageVersionStagePath+"/state", stageVersionInfo.Name+","+StageStateSuccess)
	} else if stageVersionState.RunResult == StageStateFailed {
		etcd.EtcdSet(stageVersionStagePath+"/state", stageVersionInfo.Name+","+StageStateFailed)
	}

	// notify stages dependence on current stage
	for _, v := range stageVersionInfo.To {
		if v != "" {
			bootChan <- stageVersionStagePath[:strings.LastIndex(stageVersionStagePath, "/")] + "/" + v
		}
	}
	finishChan <- stageVersionState
}

func startStageInK8S(runResultChan chan models.StageVersionState, runResult models.StageVersionState) {
	sec := rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(5) + 3

	timeStr := strconv.FormatInt(sec, 10) + "s"
	timeDur, _ := time.ParseDuration(timeStr)
	time.Sleep(timeDur)

	if rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(100) < 50 {
		runResult.RunResult = StageStateSuccess
		runResult.Detail = "run success"
	} else {
		runResult.RunResult = StageStateFailed
		runResult.Detail = "not luck"
	}

	runResultChan <- runResult
}
