package etcd

import (
	"testing"
	"github.com/containerops/vessel/models"
	"time"
	"log"
)

func Test_Etcd(t *testing.T)  {
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
	stage := &models.Stage{
		Name:"etcdStage",
		Namespace:"chenzhu",
	}
	if err := GetStage(stage); err != nil {
		log.Print(err)
	}else{
		log.Println(stage)
	}
	if err := SaveStage(fulStage()); err != nil {
		log.Print(err)
	}
	if err := GetStage(stage); err != nil {
		log.Print(err)
	}else{
		log.Println(stage)
	}
	if err := ChangeStageStatus(stage); err != nil {
		log.Print(err)
	}
	str, err := GetStageStatus(stage);
	if  err != nil {
		log.Print(err)
	}else{
		log.Println(str)
	}
	if err := SetStageStatusTTL(stage,2); err != nil {
		log.Print(err)
	}else{
		<-time.After(time.Second * time.Duration(4))
		if err := GetStage(stage); err != nil {
			log.Print(err)
		}else{
			log.Println(stage)
		}
	}
}

func fulStage() *models.Stage {
	return  &models.Stage{
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
		Dependence:[]string{"stageA","stageB","stageC"},
		Status:"OK",
	}
}

func checkPipeline() {
	pipeline := &models.Pipeline{
		Name:"etcdPipeline",
		Namespace:"chenzhu",
	}
	if err := GetPipeline(pipeline); err != nil {
		log.Print(err)
	}else{
		log.Println(pipeline)
	}
	if err := SavePipeline(fulPipeline()); err != nil {
		log.Print(err)
	}
	if err := GetPipeline(pipeline); err != nil {
		log.Print(err)
	}else{
		log.Println(pipeline)
	}
	if err := ChangePipelineStatus(pipeline); err != nil {
		log.Print(err)
	}
	str, err := GetPipelineStatus(pipeline);
	if  err != nil {
		log.Print(err)
	}else{
		log.Println(str)
	}
	if err := SetPipelineStatusTTL(pipeline,2); err != nil {
		log.Print(err)
	}else{
		<-time.After(time.Second * time.Duration(4))
		if err := GetPipeline(pipeline); err != nil {
			log.Print(err)
		}else{
			log.Println(pipeline)
		}
	}
}

func fulPipeline() *models.Pipeline {
	timeStr := time.Now().Format("2016-01-02 15:04:05")
	return  &models.Pipeline{
		Name:"etcdPipeline",
		Namespace:"chenzhu",
		Stages:[]string{"stageA","stageB","stageC"},
		CreationTimestamp:timeStr,
		DeletionTimestamp:timeStr,
		TimeoutDuration:60,
		Status:"OK",
	}
}