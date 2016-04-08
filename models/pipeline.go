package models

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/context"

	"github.com/containerops/vessel/module/etcd"
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
	etcd.Set(pipelinePath+"/allstage", strings.Join(stageNames, ","))
	for _, stage := range pipeline.Stages {
		stagePath := pipelinePath + stage.Name
		etcd.Set(stagePath+"/name", stage.Name)
		etcd.Set(stagePath+"/detail", stage.Detail)
		etcd.Set(stagePath+"/from", relationMap[stage.Name][0])
		etcd.Set(stagePath+"/to", relationMap[stage.Name][1])
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

	pipelineVersionPath := fmt.Sprintf(DEFAULT_PIPELINEVERSION_ETCD_PATH, pipelineVersion.WorkspaceId, pipelineVersion.ProjectId, pipelineVersion.Id, pipelineVersion.Id)
	etcd.Set(pipelineVersionPath+"pipelineId", strconv.FormatInt(pipeline.Id, 10))

	return pipelineVersion, nil
}

// test is the given pipeline is legal ,if legal return pipeline's stage relationMap if not return error
func isPipelineLegal(pipeline *Pipeline) (map[string][]string, error) {
	stageMap := make(map[string]*Stage, 0)
	dependenceCount := make(map[string]int, 0)
	stageRelationMap := make(map[string][]string, 0)

	for _, stage := range pipeline.Stages {
		if _, ok := stageMap[stage.Name]; !ok {
			stageMap[stage.Name] = stage

			if _, exist := stageRelationMap[stage.Name]; !exist {
				stageRelationMap[stage.Name] = make([]string, 2)
				stageRelationMap[stage.Name][0] = strings.Join(stage.From, ",")
			}

			for _, in := range stage.From {
				dependenceCount[in]++
				if _, exist := stageRelationMap[in]; !exist {
					stageRelationMap[in] = make([]string, 2)
				}
				stageRelationMap[in][1] = strings.Join(append(strings.Split(stageRelationMap[in][1], ","), stage.Name), ",")
			}
		} else {
			// has a repeat stage name ,return
			return nil, errors.New("stage has repeat name")
		}
	}

	finish := 0
	for true {
		temp := 0
		for _, stage := range stageMap {
			if dependenceCount[stage.Name] == 0 {
				finish++
				for _, out := range stage.From {
					dependenceCount[out]--
				}
				dependenceCount[stage.Name] = -1
			} else if dependenceCount[stage.Name] == -1 {
				temp++
			}

			if temp == finish || finish == len(dependenceCount) {
				break
			}
		}
	}

	if finish != len(dependenceCount) {
		return nil, errors.New("given pipeline's stage can't create a DAG")
	}

	return stageRelationMap, nil
}

// start a pipelineVersion,boot the stage and return the result
func (pipelineVersion *PipelineVersion) Boot() error {
	bootChan := make(chan string, 100)
	finishChan := make(chan string, 100)
	isDone := make(chan bool, 1)
	// get all stageName list and range all stage to start all stage
	pipelinePath := fmt.Sprintf(DEFAULT_PIPELINE_ETCD_PATH, pipelineVersion.WorkspaceId, pipelineVersion.ProjectId, pipelineVersion.PipelineId)
	pipelineVersionPath := fmt.Sprintf(DEFAULT_PIPELINEVERSION_ETCD_PATH, pipelineVersion.WorkspaceId, pipelineVersion.ProjectId, pipelineVersion.Id, pipelineVersion.Id)
	stageList, _ := etcd.Get(pipelinePath + "/allstage")
	stageNames := strings.Split(stageList.Node.Value, ",")
	sumStage := len(stageNames)

	go bootStage(bootChan, finishChan)
	go isFinish(finishChan, isDone, sumStage)

	// search all stage start stage which from is ""
	for _, stageName := range stageNames {
		stagePath := pipelinePath + stageName
		stageVersionPath := pipelineVersionPath + stageName
		fromInfo, _ := etcd.Get(stagePath + "/from")
		from := fromInfo.Node.Value
		if from == "" {
			bootChan <- stageVersionPath
		}
	}

	<-isDone
	log.Println("all job is done!")
	return nil
}

// receive bootChan's message start give stage by stage path in etcd
func bootStage(bootChan chan string, finishChan chan string) {
	stagePath := <-bootChan
	log.Println("get a job to start stage in path:", stagePath)
	go startStage(stagePath, bootChan, finishChan)
}

// count finish stage num,send single when all stage is start
func isFinish(finishChan chan string, isDone chan bool, sumStage int) {
	count := 0
	for {
		if count == sumStage {
			isDone <- true
		}

		stageName := <-finishChan
		count++
		log.Printf("stage %s is done!", stageName)
	}
}

// stage stage by given name,after start give single to finishChan
func startStage(stagePath string, bootChan, finishChan chan string) {

	// get info from etcd
	pipelineVersionPath := stagePath[:strings.LastIndex(stagePath, "/")]
	pipelineVersion := pipelineVersionPath[strings.LastIndex(pipelineVersionPath, "-")+1:]
	stageName := stagePath[strings.LastIndex(stagePath, "/")+1:]
	pipelineIdInfo, _ := etcd.Get(pipelineVersionPath + "/pipelineId")
	pipelineId := pipelineIdInfo.Node.Value
	pipelinePath := pipelineVersionPath[:strings.LastIndex(pipelineVersionPath, "/")] + "/pl-" + pipelineId + "/stage/"

	// check if the dir is exist
	stateInfo, _ := etcd.Get(stagePath + "/state")
	if stateInfo.Node.Value != "" {
		return
	}
	etcd.Set(stagePath+"/state", stageName+",1")
	// get current stage from info
	fromStageNamesInfo, _ := etcd.Get(pipelinePath + stageName + "/from")
	fromStageNames := fromStageNamesInfo.Node.Value
	stageNameMap := make(map[string]string)
	for _, stageName := range strings.Split(fromStageNames, ",") {
		stageNameMap[stageName] = stageName
	}

	count := 0
	sum := len(stageNameMap)

	// check is there is some stage finish start before this stage start
	for _, stageName := range stageNameMap {
		fromStagePath := pipelineVersionPath + "/" + pipelineVersion + "/stagev-" + pipelineVersion + "/state"
		fromStageInfo, _ := etcd.Get(fromStagePath)
		if info := fromStageInfo.Node.Value; info != "" {
			if strings.Split(info, ",")[1] == "2" {
				count++
			} else if strings.Split(info, ",")[1] == "3" {
				// if from stage is failed return
				toStageInfos, _ := etcd.Get(pipelinePath + stageName + "/to")
				for _, v := range strings.Split(toStageInfos.Node.Value, ",") {
					bootChan <- v
				}
				finishChan <- stageName
				etcd.Set(stagePath+"/state", stageName+",3")
				return
			}
		}
	}

	watcher := etcd.Watch(pipelineVersionPath + "/stagev-" + pipelineVersion + "/")
	for {
		// all stage dependence is start
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
				if strings.Split(changeStageInfo, ",")[1] == "2" {
					count++
				} else if strings.Split(changeStageInfo, ",")[1] == "3" {
					// if from stage is failed return
					toStageInfos, _ := etcd.Get(pipelinePath + stageName + "/to")
					for _, v := range strings.Split(toStageInfos.Node.Value, ",") {
						bootChan <- v
					}
					finishChan <- stageName
					etcd.Set(stagePath+"/state", stageName+",3")
					return
				}
			}
		}
	}

	sec := rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(5)

	log.Println("stage:", stageName, "will start after", sec, "s...")
	timeStr := strconv.FormatInt(sec, 10) + "s"
	timeDur, _ := time.ParseDuration(timeStr)
	time.Sleep(timeDur)

	etcd.Set(stagePath+"/state", stageName+",2")
	toStageInfos, _ := etcd.Get(pipelinePath + stageName + "/to")
	for _, v := range strings.Split(toStageInfos.Node.Value, ",") {
		bootChan <- v
	}
	finishChan <- stageName
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
