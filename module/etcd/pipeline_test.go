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

func TestSavePipeline(t *testing.T) {
	pipeline := fulPipeline()
	log.Println(pipeline, SavePipeline(pipeline))
}

func TestGetPipeline(t *testing.T) {
	pipeline := easyPipeline()
	log.Println(pipeline, GetPipeline(pipeline))
}

func TestSetPipelineStatus(t *testing.T) {
	pipeline := fulPipeline()
	log.Println(pipeline, SetPipelineStatus(pipeline))
}

func TestGetPipelineStatus(t *testing.T) {
	pipeline := easyPipeline()
	str, err := GetPipelineStatus(pipeline)
	log.Println(pipeline, str, err)
}

func TestSetCreationTimestamp(t *testing.T) {
	pipeline := easyPipeline()
	log.Println(pipeline, SetCreationTimestamp(pipeline))
}

func TestGetCreationTimestamp(t *testing.T) {
	pipeline := easyPipeline()
	str, err := GetCreationTimestamp(pipeline)
	log.Println(pipeline, str, err)
}

func TestSetDeletionTimestamp(t *testing.T) {
	pipeline := easyPipeline()
	log.Println(pipeline, SetDeletionTimestamp(pipeline))
}

func TestGetDeletionTimestamp(t *testing.T) {
	pipeline := easyPipeline()
	str, err := GetDeletionTimestamp(pipeline)
	log.Println(pipeline, str, err)
}

func TestSetPipelineResult(t *testing.T) {
	result := &models.PipelineResult{
		Name:"etcdPipelineResult",
		Namespace:"chenzhu",
		WorkspaceId: 1000,
		ProjectId:2000,
		PipelineId:"aaaaaaaaaaa",
		Status:models.STATE_STARTING,
	}
	log.Println(result, SetPipelineResult(result))
}

func TestGetPipelineResult(t *testing.T) {
	result := &models.PipelineResult{
		Name:"etcdPipelineResult",
		Namespace:"chenzhu",
	}
	log.Println(result, GetPipelineResult(result))
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
	return &models.Pipeline{
		Name:"etcdPipeline",
		Namespace:"chenzhu",
		Stages:[]string{"stageA", "stageB", "stageC"},
		TimeoutDuration:60,
		Status:models.STATE_STARTING,
	}
}