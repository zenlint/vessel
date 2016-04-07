package handler

import (
	"net/http"
	"time"

	"github.com/containerops/vessel/models"
	"github.com/containerops/vessel/module/etcd"
	//"github.com/containerops/vessel/utils"

	"gopkg.in/macaron.v1"
)

type PipelinePOSTJSON struct {
	Kind       string `json:"kind"`
	ApiVersion string `json:"apiVersion"`
	MetaData   struct {
		Name        string `json:"name"`
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

	etcd.Save()
	// get json body

	//ignore workspace & project
	//create workspace & project
	var wsInfo *models.Workspace = &models.Workspace{}
	wsInfo.Id = 10000
	wsInfo.Name = "IT-Workspace"
	wsInfo.Description = "IT-Workspace"
	wsInfo.Actived = true
	wsInfo.Created = time.Now().Unix()
	wsInfo.Updated = time.Now().Unix()
	wsInfo.Memo = ""

	var pjInfo *models.Project = &models.Project{}
	pjInfo.Id = 20000
	pjInfo.WorkspaceId = 10000
	pjInfo.Name = "IT-Project"
	pjInfo.Description = "Description"
	pjInfo.Actived = true
	pjInfo.Created = time.Now().Unix()
	pjInfo.Updated = time.Now().Unix()
	pjInfo.Memo = ""

	//post create pipeline to db
	// --pipeline save db
	// --pipeline stage save db
	//init pipeline from db to etcd
	// --pipeline to etcd
	// --pipeline stage to etcd
	// run pipeline version by etcd
	// --pipeline version to etcd
	// --pipeline stage version to etcd

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

	var plInfo *models.Pipeline = &models.Pipeline{}
	plInfo.Id = 30000
	plInfo.WorkspaceId = 10000
	plInfo.ProjectId = 20000
	plInfo.Name = "IT-Pipeline"
	plInfo.SelfLink = ""
	plInfo.Created = time.Now().Unix()
	plInfo.Updated = time.Now().Unix()
	plInfo.Labels = ""
	plInfo.Annotations = ""
	plInfo.Detail = "{}"

	var sInfo1 *models.Stage = &models.Stage{}
	sInfo1.Id = 31000
	sInfo1.WorkspaceId = 10000
	sInfo1.ProjectId = 20000
	sInfo1.PipelineId = 30000
	sInfo1.Created = time.Now().Unix()
	sInfo1.Updated = time.Now().Unix()
	sInfo1.Name = "BaseStage1"
	sInfo1.Detail = "{}"
	sInfo1.Dependences = []string{}
	plInfo.Stages = append(plInfo.Stages,sInfo1)

	var sInfo2 *models.Stage = &models.Stage{}
	sInfo2.Id = 31000
	sInfo2.WorkspaceId = 10000
	sInfo2.ProjectId = 20000
	sInfo2.PipelineId = 30000
	sInfo2.Created = time.Now().Unix()
	sInfo2.Updated = time.Now().Unix()
	sInfo2.Name = "BaseStage2"
	sInfo2.Detail = "{}"
	sInfo2.Dependences = []string{}
	plInfo.Stages = append(plInfo.Stages,sInfo2)

	var sInfo3 *models.Stage = &models.Stage{}
	sInfo3.Id = 31000
	sInfo3.WorkspaceId = 10000
	sInfo3.ProjectId = 20000
	sInfo3.PipelineId = 30000
	sInfo3.Created = time.Now().Unix()
	sInfo3.Updated = time.Now().Unix()
	sInfo3.Name = "BaseStage3"
	sInfo3.Detail = "{}"
	sInfo3.Dependences = []string{}
	plInfo.Stages = append(plInfo.Stages,sInfo3)


	var sInfo4 *models.Stage = &models.Stage{}
	sInfo4.Id = 31000
	sInfo4.WorkspaceId = 10000
	sInfo4.ProjectId = 20000
	sInfo4.PipelineId = 30000
	sInfo4.Created = time.Now().Unix()
	sInfo4.Updated = time.Now().Unix()
	sInfo4.Name = "BaseStage4"
	sInfo4.Detail = "{}"
	sInfo4.Dependences = []string{"BaseStage1", "BaseStage2"}
	plInfo.Stages = append(plInfo.Stages,sInfo4)


	var sInfo5 *models.Stage = &models.Stage{}
	sInfo5.Id = 31000
	sInfo5.WorkspaceId = 10000
	sInfo5.ProjectId = 20000
	sInfo5.PipelineId = 30000
	sInfo5.Created = time.Now().Unix()
	sInfo5.Updated = time.Now().Unix()
	sInfo5.Name = "BaseStage5"
	sInfo5.Detail = "{}"
	sInfo5.Dependences = []string{"BaseStage2", "BaseStage3"}
	plInfo.Stages = append(plInfo.Stages,sInfo5)


	var sInfo6 *models.Stage = &models.Stage{}
	sInfo6.Id = 31000
	sInfo6.WorkspaceId = 10000
	sInfo6.ProjectId = 20000
	sInfo6.PipelineId = 30000
	sInfo6.Created = time.Now().Unix()
	sInfo6.Updated = time.Now().Unix()
	sInfo6.Name = "BaseStage6"
	sInfo6.Detail = "{}"
	sInfo6.Dependences = []string{"BaseStage4", "BaseStage5"}
	plInfo.Stages = append(plInfo.Stages,sInfo6)


	var sInfo7 *models.Stage = &models.Stage{}
	sInfo7.Id = 31000
	sInfo7.WorkspaceId = 10000
	sInfo7.ProjectId = 20000
	sInfo7.PipelineId = 30000
	sInfo7.Created = time.Now().Unix()
	sInfo7.Updated = time.Now().Unix()
	sInfo7.Name = "BaseStage7"
	sInfo7.Detail = "{}"
	sInfo7.Dependences = []string{"BaseStage2", "BaseStage6"}
	plInfo.Stages = append(plInfo.Stages,sInfo7)

	var sInfo8 *models.Stage = &models.Stage{}
	sInfo8.Id = 31000
	sInfo8.WorkspaceId = 10000
	sInfo8.ProjectId = 20000
	sInfo8.PipelineId = 30000
	sInfo8.Created = time.Now().Unix()
	sInfo8.Updated = time.Now().Unix()
	sInfo8.Name = "BaseStage8"
	sInfo8.Detail = "{}"
	sInfo8.Dependences = []string{"BaseStage4", "BaseStage7"}
	plInfo.Stages = append(plInfo.Stages,sInfo8)


	// //error 9
	// var sInfo9 *models.Stage = &models.Stage{}
	// sInfo9.Id = 31000
	// sInfo9.WorkspaceId = 10000
	// sInfo9.ProjectId = 20000
	// sInfo9.PipelineId = 30000
	// sInfo9.Created = time.Now().Unix()
	// sInfo9.Updated = time.Now().Unix()
	// sInfo9.Name = "BaseStage9"
	// sInfo9.Detail = "{}"
	// sInfo9.Dependences = []string{"BaseStage4","BaseStage7","BaseStage7","BaseStage9"}
	// plInfo.Stages = append(plInfo.Stages,sInfo9)

	// var plvInfo *models.PipelineVersion = &models.PipelineVersion{}
	// plvInfo.Id = 40000
	// plvInfo.WorkspaceId = 10000
	// plvInfo.ProjectId = 20000
	// plvInfo.PipelineId = 30000
	// plvInfo.Namespace = fmt.Sprintf("%d-%d", 30000, time.Now().Unix())
	// plvInfo.SelfLink = ""
	// plvInfo.Created = time.Now().Unix()
	// plvInfo.Updated = time.Now().Unix()
	// plvInfo.Labels = ""
	// plvInfo.Annotations = ""
	// plvInfo.Detail = "{}"
	// plvInfo.Log = ""
	// plvInfo.Status = 0

	//reqStr, _ := ctx.Req.Body().String()
	// create new pipeline
	//pipeline := new(models.Pipeline)

	//ignore workspace & project
	//projectInfo := module.GetProjectInfoByName(ctx.Params(":project"))
	//pipeline.WorkspaceId = projectInfo.WorkspaceId
	//pipeline.ProjectId = projectInfo.Id
	/*

		pipeline.Name = reqData.MetaData.Name
		pipeline.SelfLink = ""
		pipeline.Labels = reqData.MetaData.Labels
		pipeline.Annotations = reqData.MetaData.Annotations
		pipeline.Created = time.Now().Unix()
		pipeline.Updated = time.Now().Unix()
		pipeline.Detail = reqStr

		module.CreatePipeline(pipeline)
	*/
	// gen new stage & point
	// verification json format,return stage & point map
	//isLegal, reason, dependenceMap := utils.GenerateDependenceMap(reqStr)
	//if !isLegal {
	//	return http.StatusOK, []byte(reason)
	//}

	//for name, define := range dependenceMap {
	/*
		point := new(models.Point)
		point.PipelineId = pipeline.Id
		point.Created = time.Now().Unix()
		point.Updated = time.Now().Unix()
		point.Name = name
		point.From = define[0]
		point.To = define[1]

		module.CreatePoint(point)

		stage := new(models.Stage)
		stage.PipelineId = pipeline.Id
		stage.Created = time.Now().Unix()
		stage.Updated = time.Now().Unix()
		stage.Name = name
		stage.Detail = define[2]

		module.CreateStage(stage)
	*/
	//}

	return http.StatusOK, []byte("ok")
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
