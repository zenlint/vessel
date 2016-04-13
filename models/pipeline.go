package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/context"

	log "github.com/Sirupsen/logrus"
)

const (
	PIPESUCCESS = iota
	PIPEERROR
	PIPERUNNING
	PIPEPENDDING
	PIPEEXCEPT
)

var (
	DEFAULT_PIPELINE_ETCD_PATH        = "/containerops/vessel/ws-%d/pj-%d/pl-%d/stage/"
	DEFAULT_PIPELINEVERSION_ETCD_PATH = "/containerops/vessel/ws-%d/pj-%d/plv-%d/stagev-%d/"
)

// type Pipeline struct {
// 	Id          int64     `json:"id"`                                        //
// 	ProjectId   int64     `json:"projectId"`                                 //
// 	UUID        string    `json:"uuid" orm:"varchar(255)"`                   //
// 	Name        string    `json:"name" orm:"varchar(255)"`                   //
// 	Description string    `json:"description" orm:"null;type(text)"`         //
// 	Actived     bool      `json:"actived" orm:"null;default(0)"`             //
// 	RootId      int64     `json:"rootId" orm:"default(0)"`                   //
// 	ParentId    int64     `json:"parentId" orm:"default(0)"`                 //
// 	Version     bool      `json:"version" orm:"default(0)"`                  //
// 	Created     time.Time `json:"created" orm:"auto_now_add;type(datetime)"` //
// 	Updated     time.Time `json:"updated" orm:"auto_now;type(datetime)"`     //
// 	Memo        string    `json:"memo" orm:"null;type(text)"`                //
// }

type Pipeline struct {
	Id          int64  `json:"id"`
	WorkspaceId int64  `json:"workspaceId"`
	ProjectId   int64  `json:"projectId"`
	Name        string `json:"name" gorm:"type:varchar(255)"`
	SelfLink    string `json:"selfLink" gorm:"type:varchar(255)"`
	Created     int64  `json:"created"`
	Updated     int64  `json:"updated"`
	Labels      string `json:"labels"`
	Annotations string `json:"annotations"`
	Detail      string `json:"detail" gorm:"type:text"`
	Stages      []*Stage
	MetaData    string `json:"metadata"`
	Spec        string `json:"spec"`
}

type PipelineVersion struct {
	Id            int64    `json:"id"`
	WorkspaceId   int64    `json:"workspaceId"`
	ProjectId     int64    `json:"projectId"`
	PipelineId    int64    `json:"pipelineId"`
	Namespace     string   `json:"namespace"`
	SelfLink      string   `json:"selfLink" gorm:"type:varchar(255)"`
	Created       int64    `json:"created"`
	Updated       int64    `json:"updated"`
	Labels        string   `json:"labels"`
	Annotations   string   `json:"annotations"`
	Detail        string   `json:"detail" gorm:"type:text"`
	StageVersions []string `json:"stageVersions"`
	Log           string   `json:"log" gorm:"type:text"`
	Status        int64    `json:"state"` // 0 not start    1 working    2 success     3 failed
	MetaData      string   `json:"metadata"`
	Spec          string   `json:"spec"`
}

// type Status struct {
// 	Id         int64     `json:"id"`                                        //
// 	PipelineId int64     `json:"pipelineId"`                                //
// 	UUID       string    `json:"uuid" orm:"varchar(255)"`                   //
// 	Resource   string    `json:"resrouce" orm:"null;type(text)"`            //
// 	ActionId   string    `json:"actionId" orm:"unique;varchar(255)"`        // Point.UUID; Stage.UUID
// 	Started    time.Time `json:"started" orm:"type(datetime)"`              //
// 	Ended      time.Time `json:"ended" orm:"type(datetime)"`                //
// 	Log        string    `json:"log" orm:"type(text)"`                      //
// 	Result     int64     `json:"result" orm:"null"`                         // Success: 0; Failure: 1
// 	Actived    bool      `json:"actived" orm:"null;default(0)"`             //
// 	Created    time.Time `json:"created" orm:"auto_now_add;type(datetime)"` //
// 	Updated    time.Time `json:"updated" orm:"auto_now;type(datetime)"`     //
// 	Memo       string    `json:"memo" orm:"null;type(text)"`                //
// }

//Create Pipeline
func (pipe *Pipeline) Create(projectId int64, name string) (error, int64) {
	return nil, 0
}

// //Create status records with same uuid
// func (pipe *Pipeline) Run(pipelineId int64) (error, string) {
// 	return nil, ""
// }

// init a pipeline generate the stage dependences, init a pipelineVersion by current pipeline and return it
func (pipeline *Pipeline) Run() (*PipelineVersion, error) {
	// first test is pipeline legal if not return err
	relationMap, err := isPipelineLegal(pipeline)
	if err != nil {
		return nil, err
	}

	stageNames := make([]string, 0)
	for _, stage := range pipeline.Stages {
		stageNames = append(stageNames, stage.Name)
	}

	pipelinePath := fmt.Sprintf(DEFAULT_PIPELINE_ETCD_PATH, pipeline.WorkspaceId, pipeline.ProjectId, pipeline.Id)
	// create ETCD dir
	EtcdSet(pipelinePath+"/allstage", strings.Join(stageNames, ","))
	for _, stage := range pipeline.Stages {
		stagePath := pipelinePath + stage.Name
		EtcdSet(stagePath+"/id", strconv.FormatInt(stage.Id, 10))
		EtcdSet(stagePath+"/name", stage.Name)
		EtcdSet(stagePath+"/detail", stage.Detail)
		EtcdSet(stagePath+"/from", relationMap[stage.Name][0])
		EtcdSet(stagePath+"/to", relationMap[stage.Name][1])
	}

	pipelineVersion := new(PipelineVersion)
	pipelineVersion.Id = time.Now().UnixNano()
	pipelineVersion.WorkspaceId = pipeline.WorkspaceId
	pipelineVersion.ProjectId = pipeline.ProjectId
	pipelineVersion.PipelineId = pipeline.Id
	pipelineVersion.Namespace = "plv" + "-" + strconv.FormatInt(pipelineVersion.Id, 10)
	pipelineVersion.SelfLink = ""
	pipelineVersion.Created = time.Now().Unix()
	pipelineVersion.Updated = time.Now().Unix()
	pipelineVersion.Labels = pipeline.Labels
	pipelineVersion.Annotations = pipeline.Annotations
	pipelineVersion.Detail = pipeline.Detail
	pipelineVersion.StageVersions = []string{strconv.FormatInt(pipelineVersion.Id, 10)}
	pipelineVersion.Status = 0

	stageVersionPath := fmt.Sprintf(DEFAULT_PIPELINEVERSION_ETCD_PATH, pipelineVersion.WorkspaceId, pipelineVersion.ProjectId, pipelineVersion.Id, pipelineVersion.Id)
	stageVersionPath = stageVersionPath[:len(stageVersionPath)-1]
	EtcdSet(stageVersionPath[:strings.LastIndex(stageVersionPath, "/")]+"/pipelineId", strconv.FormatInt(pipeline.Id, 10))

	return pipelineVersion, nil
}

// test is the given pipeline is legal ,if legal return pipeline's stage relationMap if not return error
func isPipelineLegal(pipeline *Pipeline) (map[string][]string, error) {
	stageMap := make(map[string]*Stage, 0)
	dependenceCount := make(map[string]int, 0)
	stageRelationMap := make(map[string][]string, 0)

	// regist all stage,and check repeat/nil stage name
	for _, stage := range pipeline.Stages {
		if stage.Name == "" {
			return nil, errors.New("stage has a nil name")
		}
		if _, exist := stageMap[stage.Name]; !exist {
			stageMap[stage.Name] = stage
		} else {
			// has a repeat stage name ,return
			return nil, errors.New("stage has repeat name:" + stage.Name)
		}
	}

	// init stage dependence count
	for stageName, _ := range stageMap {
		dependenceCount[stageName] = 0
	}

	// count stage dependence
	for _, stage := range stageMap {
		for _, from := range stage.From {
			dependenceCount[from]++
		}
	}

	// check DAG
	//if AnnulusTag == nowReleaseStageCount or nowReleaseStageCount == len(dependenceCount) then exit for,if nowReleaseStageCount == len(dependenceCount) then isDAG,else isNotDAG
	nowReleaseStageCount := 0
	for true {

		annulusTag := 0
		for stageName, stage := range stageMap {
			if dependenceCount[stageName] == 0 {
				nowReleaseStageCount++
				for _, from := range stage.From {
					dependenceCount[from]--
				}

				dependenceCount[stage.Name] = -1
			} else if dependenceCount[stageName] == -1 {
				annulusTag++
			}
		}

		if annulusTag == nowReleaseStageCount || nowReleaseStageCount == len(dependenceCount) {
			break
		}
	}

	if nowReleaseStageCount != len(dependenceCount) {
		return nil, errors.New("given pipeline's stage can't create a DAG")
	}

	// generate stage relationMap
	// stageRelationMap := map[stageName]{"stage.From","stage.To"}
	for stageName, stage := range stageMap {
		if _, exist := stageRelationMap[stageName]; !exist {
			stageRelationMap[stageName] = make([]string, 2)
		}
		stageRelationMap[stageName][0] = strings.Join(stage.From, ",")

		for _, from := range stage.From {
			if _, exist := stageRelationMap[from]; !exist {
				stageRelationMap[from] = make([]string, 2)
			}
			stageRelationMap[from][1] = strings.Join(append(strings.Split(stageRelationMap[from][1], ","), stageName), ",")
		}
	}

	return stageRelationMap, nil

}

const (
	StageStateNotStart = "not start"
	StageStateStarting = "working"
	StageStateSuccess  = "success"
	StageStateFailed   = "failed"
)

type StageVersionState struct {
	PipelineId        string `json:"pipelineId"`
	PipelineVersionId string `json:"pipelineVersionId"`
	StageId           string `json:"pipelineVersionId"`
	StageVersionId    string `json:"stageVersionId"`
	StageName         string `json:"stageName"`
	RunResult         string `json:"runResult"`
	Detail            string `json:"detail"`
}

// start a pipelineVersion,boot the stage and return the result
func (pipelineVersion *PipelineVersion) Boot() string {
	bootChan := make(chan string, 100)
	finishChan := make(chan StageVersionState, 100)
	stageVersionStateChan := make(chan string, 1)
	notifyBootDone := make(chan bool, 1)
	// get all stageName list and range all stage to start all stage
	pipelinePath := fmt.Sprintf(DEFAULT_PIPELINE_ETCD_PATH, pipelineVersion.WorkspaceId, pipelineVersion.ProjectId, pipelineVersion.PipelineId)
	pipelineVersionPath := fmt.Sprintf(DEFAULT_PIPELINEVERSION_ETCD_PATH, pipelineVersion.WorkspaceId, pipelineVersion.ProjectId, pipelineVersion.Id, pipelineVersion.Id)
	stageList, _ := EtcdGet(pipelinePath + "/allstage")
	stageNames := strings.Split(stageList.Node.Value, ",")
	sumStage := len(stageNames)

	go bootStage(bootChan, finishChan, notifyBootDone)
	go isFinish(finishChan, stageVersionStateChan, sumStage, notifyBootDone)

	// search all stage, start stage which from is ""
	for _, stageName := range stageNames {
		if stageName != "" {
			stagePath := pipelinePath + stageName
			stageVersionPath := pipelineVersionPath + stageName
			fromInfo, _ := EtcdGet(stagePath + "/from")
			from := fromInfo.Node.Value
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
func bootStage(bootChan chan string, finishChan chan StageVersionState, notifyBootDone chan bool) {
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
func isFinish(finishChan chan StageVersionState, stageVersionStateChan chan string, sumStage int, notifyBootDone chan bool) {
	finishStageNum := 0
	stageVersionStateList := make([]StageVersionState, 0)

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
func startStage(stageVersionStagePath string, bootChan chan string, finishChan chan StageVersionState) {
	// get info from etcd
	// stageName
	stageName := stageVersionStagePath[strings.LastIndex(stageVersionStagePath, "/")+1:]
	// stageVersionPath
	stageVersionPath := stageVersionStagePath[:strings.LastIndex(stageVersionStagePath, "/")]
	// pipelineVersionPath
	pipelineVersionPath := stageVersionPath[:strings.LastIndex(stageVersionPath, "/")]
	// pipelineVersionId
	pipelineVersionId := pipelineVersionPath[strings.LastIndex(pipelineVersionPath, "-")+1:]
	// pipelineId
	pipelineIdInfo, _ := EtcdGet(pipelineVersionPath + "/pipelineId")
	pipelineId := pipelineIdInfo.Node.Value
	// stagePath
	stagePath := pipelineVersionPath[:strings.LastIndex(pipelineVersionPath, "/")] + "/pl-" + pipelineId + "/stage"
	// stageId
	stageIdInfo, _ := EtcdGet(stagePath + "/" + stageName + "/id")
	stageId := stageIdInfo.Node.Value
	// pipelinePath
	// pipelinePath := stagePath[:strings.LastIndex(stagePath, "/")]

	// check if the dir is exist
	stateInfo, _ := EtcdGet(stageVersionPath + "/" + stageName + "/state")
	if stateInfo != nil && stateInfo.Node.Value != "" {
		return
	}

	EtcdSet(stageVersionPath+"/state", stageName+",1")
	// get current stage from info
	fromStageNamesInfo, _ := EtcdGet(stagePath + "/" + stageName + "/from")

	fromStageNames := fromStageNamesInfo.Node.Value
	stageNameMap := make(map[string]string)
	for _, stageName := range strings.Split(fromStageNames, ",") {
		if stageName != "" {
			stageNameMap[stageName] = stageName
		}
	}

	count := 0
	sum := len(stageNameMap)
	var stageVersionState StageVersionState
	stageVersionState.PipelineId = pipelineId
	stageVersionState.PipelineVersionId = pipelineVersionId
	stageVersionState.StageId = stageId
	stageVersionState.StageVersionId = pipelineVersionId
	stageVersionState.StageName = stageName

	// check is there is some stage finish start before this stage start
	for _, fromStageName := range stageNameMap {
		fromStageVersionStatePath := stageVersionPath + "/" + fromStageName + "/state"
		fromStageVersionStateInfo, _ := EtcdGet(fromStageVersionStatePath)
		if fromStageVersionStateInfo != nil {
			info := fromStageVersionStateInfo.Node.Value
			if strings.Split(info, ",")[1] == StageStateSuccess {
				count++
			} else if strings.Split(info, ",")[1] == StageStateFailed {
				// if from stage is failed return
				toStageInfos, _ := EtcdGet(stagePath + "/" + stageName + "/to")
				if toStageInfos != nil {
					for _, v := range strings.Split(toStageInfos.Node.Value, ",") {
						if v != "" {
							bootChan <- stageVersionPath + "/" + v
						}
					}
				}
				stageVersionState.RunResult = StageStateFailed
				stageVersionState.Detail = "pre stage " + strings.Split(info, ",")[0] + " is failed"
				finishChan <- stageVersionState
				EtcdSet(stageVersionStagePath+"/state", stageName+","+StageStateFailed)
				return
			}
		}
	}

	watcher := EtcdWatch(stageVersionPath + "/")
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
					// if from stage is failed return
					toStageInfos, _ := EtcdGet(stagePath + "/" + stageName + "/to")
					if toStageInfos != nil {
						for _, v := range strings.Split(toStageInfos.Node.Value, ",") {
							if v != "" {
								bootChan <- stageVersionPath + "/" + v
							}
						}
					}
					stageVersionState.RunResult = StageStateFailed
					stageVersionState.Detail = "pre stage " + strings.Split(changeStageInfo, ",")[0] + " is failed"
					finishChan <- stageVersionState
					EtcdSet(stageVersionStagePath+"/state", stageName+","+StageStateFailed)
					return
				}
			}
		}
	}

	// start run stage
	stageStartFinish := make(chan StageVersionState, 1)
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
		EtcdSet(stageVersionStagePath+"/state", stageName+","+StageStateSuccess)
	} else if stageVersionState.RunResult == StageStateFailed {
		EtcdSet(stageVersionStagePath+"/state", stageName+","+StageStateFailed)
	}

	// notify stages dependence on current stage
	toStageInfos, _ := EtcdGet(stagePath + "/" + stageName + "/to")
	for _, v := range strings.Split(toStageInfos.Node.Value, ",") {
		if v != "" {
			bootChan <- stageVersionPath + "/" + v
		}
	}
	finishChan <- stageVersionState
}

func startStageInK8S(runResultChan chan StageVersionState, runResult StageVersionState) {
	// runResult := <-runResultChan
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

//List all run history with pipelineId
func (pipe *Pipeline) History(pipelineId int64) (error, []string) {
	return nil, []string{}
}

//Return all point and stage status with status uuid
func (pipe *Pipeline) Status(uuid string) (error, []int64) {
	return nil, []int64{}
}

//Stop run
func (pipe *Pipeline) Stop(uuid string) error {
	return nil
}

//Clean run resources
func (pipe *Pipeline) Clean(uuid string) error {
	return nil
}

func (pipe *Pipeline) Copy(uuid string) (error, string) {
	return nil, ""
}
