package handler

import (
	"net/http"

	"gopkg.in/macaron.v1"
)

type WorkspacePOSTJSON struct {
}

func V1POSTWorkspaceHandler(ctx *macaron.Context, workspace WorkspacePOSTJSON) (int, []byte) {
	return http.StatusOK, []byte("")
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
