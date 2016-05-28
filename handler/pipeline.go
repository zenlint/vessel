package handler

import (
	// "net/http"
	"gopkg.in/macaron.v1"
	"github.com/containerops/vessel/models"
)

func V1POSTPipelineHandler(ctx *macaron.Context, reqData models.PipelineSpecTemplate) (int, []byte) {
	return 0, []byte("")
}

func V1PUTPipelineHandler(ctx *macaron.Context) (int, []byte) {
	return 0, []byte("")
}

func V1GETPipelineHandler(ctx *macaron.Context) (int, []byte) {
	return 0, []byte("")
}

func V1DELETEPipelineHandler(ctx *macaron.Context, reqData models.PipelineSpecTemplate) (int, []byte) {	
	return 0, []byte("")
}

func V1RunPipelineHandler(ctx *macaron.Context) (int, []byte) {
	return 0, []byte("")
}

func V1StatusGETHandler(ctx *macaron.Context) (int, []byte) {
	return 0, []byte("")
}
