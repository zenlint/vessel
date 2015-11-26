package handler

import (
	"net/http"

	"gopkg.in/macaron.v1"
)

type StagePOSTJSON struct {
}

func V1POSTStageHandler(ctx *macaron.Context, stage StagePOSTJSON) (int, []byte) {
	return http.StatusOK, []byte("")
}

func V1PUTStageHandler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

func V1GETStageHandler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

func V1DELETEStageHandler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}
