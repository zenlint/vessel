package etcd

import (
	"testing"
	"log"
	"time"

	"github.com/containerops/vessel/models"
)

func init() {
	clientEtcd()
}

func TestSaveStage(t *testing.T) {
	stage := fulStage()
	log.Println(stage, SaveStage(stage))
}

func TestGetStage(t *testing.T) {
	stage := easyStage()
	log.Println(stage, GetStage(stage))
}

func TestSetStageStatus(t *testing.T) {
	stage := easyStage()
	stage.Status = models.STATE_DELETED
	log.Println(stage, SetStageStatus(stage))
}

func TestGetStageStatus(t *testing.T) {
	stage := easyStage()
	str, err := GetStageStatus(stage)
	log.Println(stage, str, err)
}

func TestSetStageResult(t *testing.T) {
	result := &models.StageResult{
		Namespace:"etcdStageResult",
		Id:"bbbbbbbbbb",
		Name:"stageNamea",
		Result:models.RESULT_SUCCESS,
		Detail:"VVVVVVVV",
	}
	log.Println(SetStageResult(result))
}

func TestGetStageResult(t *testing.T) {
	result := &models.StageResult{
		Namespace:"etcdStageResult",
		Name:"stageNamea",
	}
	log.Println(result,GetStageResult(result))
}

func TestSetStageTTL(t *testing.T) {
	stage := easyStage()
	log.Println(stage, SetStageTTL(stage, 2))

	<-time.After(time.Second * time.Duration(4))
	log.Println(stage, GetStage(stage))
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
		StatusCheckUrl:"/heath",
		StatusCheckInterval:30,
		StatusCheckCount:3,
		EnvName:"",
		EnvValue:"",
		Dependence:"stageA,stageB,stageC",
		Status:models.STATE_SUCCESS,
	}
}