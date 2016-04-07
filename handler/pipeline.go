package handler

import (
	"net/http"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/macaron.v1"

	"github.com/containerops/vessel/models"
	// "github.com/containerops/vessel/module/etcd"
)

type PipelinePOSTJSON struct {
	Kind       string `json:"kind"`
	ApiVersion string `json:"apiVersion"`
	MetaData   struct {
		Name        string `json:"name"`
		Workspace   string `json:"workspace"`
		Project     string `json:"project"`
		Namespace   string `json:"namespace"`
		SelfLink    string `json:"selfLink"`
		Labels      string `json:"labels"`
		Annotations string `json:"annotations"`
	} `json:"metadata"`
	Spec []struct {
		Name                string `json:"name"`
		Dependence          string `json:"dependence"`
		Kind                string `json:"kind"`
		StatusCheckUrl      string `json:"statusCheckUrl"`
		StatusCheckInterval int64  `json:"statusCheckInterval"`
		StatusCheckCount    int64  `json:"statusCheckCount"`
	} `json:"spec"`
}

func V1POSTPipelineHandler(ctx *macaron.Context, reqData PipelinePOSTJSON) (int, []byte) {
	/*
		etcd path /vessel/ws-xxx/pj-xxx/pl-xxx/plv-xxx/stage-xxx/

		plv-xxx  -> k8s namespace

		demo:
		/vessel/ws-xxx/pj-xxx/pl-xxx1/stage/stage1/...
		/vessel/ws-xxx/pj-xxx/pl-xxx1/stage/stage2...
		/vessel/ws-xxx/pj-xxx/pl-xxx2/stage/stage1/...
		/vessel/ws-xxx/pj-xxx/pl-xxx2/stage/stage2...
		/vessel/ws-xxx/pj-xxx/pl-xxx1/plv-xxx/stagev-xxx/name
		/vessel/ws-xxx/pj-xxx/pl-xxx1/plv-xxx/stagev-xxx/dependence/Dependence1ServicesName <--need watch
		/vessel/ws-xxx/pj-xxx/pl-xxx1/plv-xxx/stagev-xxx/dependence/Dependence2ServicesName <--need watch
		/vessel/ws-xxx/pj-xxx/pl-xxx1/plv-xxx/stagev-xxx/check/check_status_url
		/vessel/ws-xxx/pj-xxx/pl-xxx1/plv-xxx/stagev-xxx/check/check_status_interval
		/vessel/ws-xxx/pj-xxx/pl-xxx1/plv-xxx/stagev-xxx/check/check_status_count
	*/
	createPipelineAndStage(reqData)
	return http.StatusOK, []byte("ok")
}

func createPipelineAndStage(plJson PipelinePOSTJSON) (*models.Pipeline, error) {
	var plInfo *models.Pipeline = &models.Pipeline{}
	plInfo.Id = time.Now().UnixNano()
	//ignore workspace & project
	plInfo.WorkspaceId = 10000
	plInfo.ProjectId = 20000
	plInfo.Name = plJson.MetaData.Name
	plInfo.Created = time.Now().Unix()
	plInfo.Updated = time.Now().Unix()
	plInfo.Labels = plJson.MetaData.Labels
	plInfo.Annotations = plJson.MetaData.Annotations
	//ignore plJson Detail
	plInfo.Detail = ""

	for _, value := range plJson.Spec {
		var sInfo *models.Stage = &models.Stage{}
		sInfo.Id = time.Now().Unix()
		sInfo.WorkspaceId = plInfo.WorkspaceId
		sInfo.ProjectId = plInfo.ProjectId
		sInfo.PipelineId = plInfo.Id
		sInfo.Created = time.Now().Unix()
		sInfo.Updated = time.Now().Unix()
		sInfo.Name = value.Name
		//ignore Stage Detail
		//StatusCheckUrl to Detail
		//StatusCheckInterval to Detail
		//StatusCheckCount to Detail
		sInfo.Detail = ""
		sInfo.Dependences = strings.Split(value.Dependence, ",")
		plInfo.Stages = append(plInfo.Stages, sInfo)
	}

	log.Error(plInfo)
	for _, value := range plInfo.Stages {
		log.Error(value)
	}

	return plInfo, nil
}

func V1PUTPipelineHandler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

func V1GETPipelineHandler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

func V1DELETEPipelineHandler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

func V1RunPipelineHandler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

func V1StatusGETHandler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}
