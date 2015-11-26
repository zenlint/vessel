package handler

import (
	"net/http"

	"gopkg.in/macaron.v1"
)

type PipelinePOSTJSON struct {
}

func V1POSTPipelineHandler(ctx *macaron.Context, pipeline PipelinePOSTJSON) (int, []byte) {
	return http.StatusOK, []byte("")
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
