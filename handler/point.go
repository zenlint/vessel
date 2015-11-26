package handler

import (
	"net/http"

	"gopkg.in/macaron.v1"
)

type PointPOSTJSON struct {
}

func V1POSTPointHandler(ctx *macaron.Context, point PointPOSTJSON) (int, []byte) {
	return http.StatusOK, []byte("")
}

func V1PUTPointHandler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

func V1GETPointHandler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

func V1DELETEPointHandler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}
