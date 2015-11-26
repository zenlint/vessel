package handler

import (
	"net/http"

	"gopkg.in/macaron.v1"
)

type ProjectPOSTJSON struct {
}

func V1POSTProjectHandler(ctx *macaron.Context, project ProjectPOSTJSON) (int, []byte) {
	return http.StatusOK, []byte("")
}

func V1PUTProjectHandler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

func V1GETProjectHandler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

func V1DELETEProjectHandler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}
