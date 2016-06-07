package handler

import (
	"github.com/containerops/vessel/models"
	"gopkg.in/macaron.v1"
)

// V1POSTPipelineHandler pipeline handler for vessel v1 version by HTTP POST
func V1POSTPipelineHandler(ctx *macaron.Context, reqData models.PipelineSpecTemplate) (int, []byte) {
	return 0, []byte("")
}

// V1PUTPipelineHandler pipeline handler for vessel v1 version by HTTP PUT
func V1PUTPipelineHandler(ctx *macaron.Context) (int, []byte) {
	return 0, []byte("")
}

// V1GETPipelineHandler pipeline handler for vessel v1 version by HTTP GET
func V1GETPipelineHandler(ctx *macaron.Context) (int, []byte) {
	return 0, []byte("")
}

// V1DELETEPipelineHandler pipeline handler for vessel v1 version by HTTP DELETE
func V1DELETEPipelineHandler(ctx *macaron.Context, reqData models.PipelineSpecTemplate) (int, []byte) {
	return 0, []byte("")
}
