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

func TestGetStage(t *testing.T) {
	stage := easyStage()
	log.Println(stage, GetStage(stage))
}

func TestSaveStage(t *testing.T) {
	stage := fulStage()
	log.Println(stage, SaveStage(stage))

	stage = easyStage()
	log.Println(stage, GetStage(stage))
}

func TestGetStageStatus(t *testing.T) {
	stage := easyStage()

	str, err := GetStageStatus(stage)
	log.Println(stage, str, err)
}

func TestChangeStageStatus(t *testing.T) {
	stage := easyStage()
	stage.Status = models.StateSuccess
	log.Println(stage, ChangeStageStatus(stage))

	stage = easyStage()
	str, err := GetStageStatus(stage)
	log.Println(stage, str, err)
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
		StatusCheckLink:"/heath",
		StatusCheckInterval:30,
		StatusCheckCount:3,
		EnvName:"",
		EnvValue:"",
		Dependence:[]string{"stageA", "stageB", "stageC"},
		Status:models.StateStarting,
	}
}