package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	//"github.com/ngaut/log"
	"gopkg.in/macaron.v1"

	"github.com/containerops/vessel/models"
)

type ProjectPOSTJSON struct {
	Name        string `from:"name" binding:"Required"`
	Description string `from:"description" binding:"Required"`
}

func V1POSTProjectHandler(ctx *macaron.Context, pj ProjectPOSTJSON) (int, []byte) {
	p := models.Project{}
	wid, _ := strconv.ParseInt(ctx.Params(""), 0, 64)

	if id, err := p.Create(wid, pj.Name, pj.Description); err != nil {
		//	log.Errorf("[vessel] Create project error: %s", err.Error())

		result, _ := json.Marshal(map[string]string{"status": "Error", "message": err.Error()})
		return http.StatusBadRequest, result
	} else {
		//	log.Infof("[vessel] Create workspace successfully, id is %d", id)

		result, _ := json.Marshal(map[string]int64{"id": id})

		return http.StatusOK, result
	}
}

type ProjectPUTJSON struct {
	Name        string `from:"name" binding:"Required"`
	Description string `from:"description" binding:"Required"`
}

func V1PUTProjectHandler(ctx *macaron.Context, pj ProjectPUTJSON) (int, []byte) {
	p := models.Project{}
	pid, _ := strconv.ParseInt(ctx.Params(":project"), 0, 64)

	if err := p.Put(pid, pj.Name, pj.Description); err != nil {
		//		log.Errorf("[vessel] Put project %d error: %s", pid, err.Error())

		result, _ := json.Marshal(map[string]string{"status": "Error", "message": err.Error()})
		return http.StatusBadRequest, result
	} else {
		//		log.Infof("[vessel] Put project %d successfully.", pid)

		result, _ := json.Marshal(map[string]int64{"id": pid})
		return http.StatusOK, result
	}
}

func V1GETProjectHandler(ctx *macaron.Context) (int, []byte) {
	p := models.Project{}
	pid, _ := strconv.ParseInt(ctx.Params(":project"), 0, 64)

	if pj, err := p.Get(pid); err != nil {
		//		log.Errorf("[vessel] Get project error: %s", err.Error())

		result, _ := json.Marshal(map[string]string{"status": "Error", "message": err.Error()})
		return http.StatusBadRequest, result
	} else {
		//		log.Infof("[vessel] Get project data successfully: %d", pid)

		result, _ := json.Marshal(pj)
		return http.StatusOK, result
	}
}

func V1DELETEProjectHandler(ctx *macaron.Context) (int, []byte) {
	p := models.Workspace{}
	pid, _ := strconv.ParseInt(ctx.Params(":project"), 0, 64)

	if err := p.Delete(pid); err != nil {
		//		log.Errorf("[vessel] Delete project error: %s", err.Error())

		result, _ := json.Marshal(map[string]string{"status": "Error", "message": err.Error()})
		return http.StatusBadRequest, result
	} else {
		//		log.Infof("[vessel] Delete project data successfully: %d", pid)

		result, _ := json.Marshal(map[string]int64{"id": pid})
		return http.StatusOK, result
	}
}
