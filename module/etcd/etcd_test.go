package etcd

import (
	"testing"
	"github.com/containerops/vessel/models"
	"time"
	"log"
)

func Test_Etcd(t *testing.T) {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	settingPoints := []map[string]string{
		map[string]string{
			"host": "127.0.0.1",
			"port": "4001",
		},
		map[string]string{
			"host": "localhost",
			"port": "4001",
		},
	}
	CreateClient(settingPoints)
	checkStage()
	checkPipeline()
}

func checkStage() {
	stage := easyStage()
	if err := GetStage(stage); err != nil {
		log.Print(err)
	} else {
		log.Println(stage)
	}

	if err := SaveStage(fulStage()); err != nil {
		log.Print(err)
	}

	stage = easyStage()
	if err := GetStage(stage); err != nil {
		log.Print(err)
	} else {
		log.Println(stage)
	}

	stage.Status = "OK"
	if err := ChangeStageStatus(stage); err != nil {
		log.Print(err)
	}

	stage = easyStage()
	str, err := GetStageStatus(stage);
	if err != nil {
		log.Print(err)
	} else {
		log.Println(str,stage)
	}

	stage.Status = "Delete"
	if err := SetStageStatusTTL(stage, 2); err != nil {
		log.Print(err)
	}

	<-time.After(time.Second * time.Duration(4))
	if err := GetStage(stage); err != nil {
		log.Print(err)
	} else {
		log.Println(stage)
	}
}

func easyStage() *models.Stage {
	return &models.Stage{
		Name:"etcdStage",
		Namespace:"chenzhu",
	}
}

func fulStage() *models.Stage {
	return &models.Stage{
		Name:"etcdStage",
		Namespace:"chenzhu",
		Replicas:3,
		Image:"unknow",
		Port:80,
		StatusCheckLink:"/heath",
		StatusCheckInterval:30,
		StatusCheckCount:3,
		EnvName:"",
		EnvValue:"",
		Dependence:[]string{"stageA", "stageB", "stageC"},
		Status:"Working",
	}
}

func checkPipeline() {
	pipeline := easyPipeline()
	if err := GetPipeline(pipeline); err != nil {
		log.Print(err)
	} else {
		log.Println(pipeline)
	}

	if err := SavePipeline(fulPipeline()); err != nil {
		log.Print(err)
	}

	pipeline = easyPipeline()
	if err := GetPipeline(pipeline); err != nil {
		log.Print(err)
	} else {
		log.Println(pipeline)
	}

	pipeline.Status = "OK"
	if err := ChangePipelineStatus(pipeline); err != nil {
		log.Print(err)
	}

	pipeline = easyPipeline()
	str, err := GetPipelineStatus(pipeline);
	if err != nil {
		log.Print(err)
	} else {
		log.Println(str, pipeline)
	}

	pipeline.Status = "Delete"
	if err := SetPipelineStatusTTL(pipeline, 2); err != nil {
		log.Print(err)
	}

	<-time.After(time.Second * time.Duration(4))
	if err := GetPipeline(pipeline); err != nil {
		log.Print(err)
	} else {
		log.Println(pipeline)
	}
}

func easyPipeline() *models.Pipeline {
	return &models.Pipeline{
		Name:"etcdPipeline",
		Namespace:"chenzhu",
	}
}

func fulPipeline() *models.Pipeline {
	timeStr := time.Now().Format("2016-01-02 15:04:05")
	return &models.Pipeline{
		Name:"etcdPipeline",
		Namespace:"chenzhu",
		Stages:[]string{"stageA", "stageB", "stageC"},
		CreationTimestamp:timeStr,
		DeletionTimestamp:timeStr,
		TimeoutDuration:60,
		Status:"Working",
	}
}