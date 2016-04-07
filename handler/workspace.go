package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	//"github.com/ngaut/log"
	"gopkg.in/macaron.v1"

	"github.com/containerops/vessel/models"
)

type WorkspacePOSTJSON struct {
	Name        string `from:"name" binding:"Required"`
	Description string `from:"description" binding:"Required"`
}

func V1POSTWorkspaceHandler(ctx *macaron.Context, ws WorkspacePOSTJSON) (int, []byte) {
	w := models.Workspace{}

	if id, err := w.Create(ws.Name, ws.Description); err != nil {
		//	log.Errorf("[vessel] Create workspace error: %s", err.Error())

		result, _ := json.Marshal(map[string]string{"status": "Error", "message": err.Error()})
		return http.StatusBadRequest, result
	} else {
		//	log.Infof("[vessel] Create workspace successfully, id is: %d", id)

		result, _ := json.Marshal(map[string]int64{"id": id})
		return http.StatusOK, result
	}
}

type WorkspacePUTJSON struct {
	Name        string `from:"name" binding:"Required"`
	Description string `from:"description" binding:"Required"`
}

func V1PUTWorkspaceHandler(ctx *macaron.Context, ws WorkspacePUTJSON) (int, []byte) {
	w := models.Workspace{}
	wid, _ := strconv.ParseInt(ctx.Params(":workspace"), 0, 64)

	if err := w.Put(wid, ws.Name, ws.Description); err != nil {
		//		log.Errorf("[vessel] Update workspace error: %s", err.Error())

		result, _ := json.Marshal(map[string]string{"status": "Error", "message": err.Error()})
		return http.StatusBadRequest, result
	} else {
		//		log.Infof("[vessel] Put workspace %d successfully", wid)

		result, _ := json.Marshal(map[string]int64{"id": wid})
		return http.StatusOK, result
	}
}

func V1GETWorkspaceHandler(ctx *macaron.Context) (int, []byte) {
	w := models.Workspace{}
	wid, _ := strconv.ParseInt(ctx.Params(":workspace"), 0, 64)

	if ws, err := w.Get(wid); err != nil {
		//		log.Errorf("[vessel] Get workspace error: %s", err.Error())

		result, _ := json.Marshal(map[string]string{"status": "Error", "message": err.Error()})
		return http.StatusBadRequest, result
	} else {
		//		log.Infof("[vessel] Get workspace data successfully: %d", wid)

		result, _ := json.Marshal(ws)
		return http.StatusOK, result
	}
}

func V1DELETEWorkspaceHandler(ctx *macaron.Context) (int, []byte) {
	w := models.Workspace{}
	wid, _ := strconv.ParseInt(ctx.Params(":workspace"), 0, 64)

	if err := w.Delete(wid); err != nil {
		//		log.Errorf("[vessel] Delete workspace error: %s", err.Error())

		result, _ := json.Marshal(map[string]string{"status": "Error", "message": err.Error()})
		return http.StatusBadRequest, result
	} else {
		//		log.Infof("[vessel] Delete workspace data successfully: %d", wid)

		result, _ := json.Marshal(map[string]int64{"id": wid})
		return http.StatusOK, result
	}
}
