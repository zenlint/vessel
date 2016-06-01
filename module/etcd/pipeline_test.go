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

func TestGetPipeline(t *testing.T) {
	pipeline := easyPipeline()
	log.Println(pipeline, GetPipeline(pipeline))
}

func TestSavePipeline(t *testing.T) {
	pipeline := fulPipeline()
	log.Println(pipeline, SavePipeline(pipeline))

	pipeline = easyPipeline()
	log.Println(pipeline, GetPipeline(pipeline))
}

func TestGetPipelineStatus(t *testing.T) {
	pipeline := easyPipeline()
	str, err := GetPipelineStatus(pipeline)
	log.Println(pipeline, str, err)
}

func TestChangePipelineStatus(t *testing.T) {
	pipeline := easyPipeline()
	pipeline.Status = models.StateSuccess
	log.Println(pipeline, ChangePipelineStatus(pipeline))

	pipeline = easyPipeline()
	str, err := GetPipelineStatus(pipeline)
	log.Println(pipeline, str, err)
}

func TestSetCreationTimestamp(t *testing.T) {
	pipeline := easyPipeline()
	pipeline.CreationTimestamp = time.Now().Format("2016-01-02 15:04:05")
	log.Println(pipeline, SetCreationTimestamp(pipeline))

	log.Println(pipeline, GetPipeline(pipeline))
}

func TestSetDeletionTimestamp(t *testing.T) {
	pipeline := easyPipeline()
	pipeline.DeletionTimestamp = time.Now().Format("2016-01-02 15:04:05")
	log.Println(pipeline, SetDeletionTimestamp(pipeline))

	log.Println(pipeline, GetPipeline(pipeline))
}

func TestSetPipelineTTL(t *testing.T) {
	pipeline := easyPipeline()
	log.Println(pipeline, SetPipelineTTL(pipeline, 2))

	<-time.After(time.Second * time.Duration(4))
	log.Println(pipeline, GetPipeline(pipeline))
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
		Status:models.StateStarting,
	}
}