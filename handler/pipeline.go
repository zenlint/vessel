package handler

import (
	// "net/http"
	"gopkg.in/macaron.v1"
	"github.com/containerops/vessel/models"
)

// V1POSTPipelineHandler
func V1POSTPipelineHandler(ctx *macaron.Context, reqData models.PipelineSpecTemplate) (int, []byte) {
	return 0, []byte("")
}

// V1PUTPipelineHandler
func V1PUTPipelineHandler(ctx *macaron.Context) (int, []byte) {
	return 0, []byte("")
}

// V1GETPipelineHandler
func V1GETPipelineHandler(ctx *macaron.Context) (int, []byte) {
	return 0, []byte("")
}

// V1DELETEPipelineHandler
func V1DELETEPipelineHandler(ctx *macaron.Context, reqData models.PipelineSpecTemplate) (int, []byte) {	
	return 0, []byte("")
}

// V1RunPipelineHandler
func V1RunPipelineHandler(ctx *macaron.Context) (int, []byte) {
	return 0, []byte("")
}

// V1StatusGETHandler
func V1StatusGETHandler(ctx *macaron.Context) (int, []byte) {
	return 0, []byte("")
}
