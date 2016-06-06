package handler

import (
	"net/http"

	"gopkg.in/macaron.v1"
)

// StagePOSTJSON
type StagePOSTJSON struct {
}

// V1POSTStageHandler
func V1POSTStageHandler(ctx *macaron.Context, stage StagePOSTJSON) (int, []byte) {
	return http.StatusOK, []byte("")
}

// V1PUTStageHandler
func V1PUTStageHandler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

// V1GETStageHandler
func V1GETStageHandler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

// V1DELETEStageHandler
func V1DELETEStageHandler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}
