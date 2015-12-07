package handler

import (
	"net/http"

	"gopkg.in/macaron.v1"
)

type ProjectPOSTJSON struct {
	Name        string `from:"name" binding:"Required"`
	Description string `from:"description" binding:"Required"`
}

func V1POSTProjectHandler(ctx *macaron.Context, project ProjectPOSTJSON) (int, []byte) {
	return http.StatusOK, []byte("")
}

type ProjectPutJSON struct {
	Description string `from:"description" binding:"Required"`
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
