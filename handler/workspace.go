package handler

import (
	"encoding/json"
	"net/http"

	"github.com/containerops/vessel/models"
	"github.com/ngaut/log"
	"gopkg.in/macaron.v1"
)

type WorkspacePOSTJSON struct {
	Name        string `from:"name" binding:"Required"`
	Description string `from:"description" binding:"Required"`
}

func V1POSTWorkspaceHandler(ctx *macaron.Context, workspace WorkspacePOSTJSON) (int, []byte) {
	w := models.Workspace{}

	if id, err := w.Create(workspace.Name, workspace.Description); err != nil {
		log.Errorf("[vessel] Create workspace error: %s", err.Error())

		result, _ := json.Marshal(map[string]string{"status": "Error", "message": err.Error()})
		return http.StatusBadRequest, result
	} else {
		log.Errorf("[vessel] Create workspace successfully, id is: %d", id)

		result, _ := json.Marshal(map[string]int64{"id": id})
		return http.StatusOK, result
	}
}

func V1PUTWorkspaceHandler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

func V1GETWorkspaceHandler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

func V1DELETEWorkspaceHandler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}
