package handler

import (
	"net/http"

	"gopkg.in/macaron.v1"
)

// StagePOSTJSON struct for json data by HTTP POST
type StagePOSTJSON struct {
}

// V1POSTStageHandler stage handler for vessel v1 version by HTTP POST
func V1POSTStageHandler(ctx *macaron.Context, stage StagePOSTJSON) (int, []byte) {
	return http.StatusOK, []byte("")
}

// V1PUTStageHandler stage handler for vessel v1 version by HTTP PUT
func V1PUTStageHandler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

// V1GETStageHandler stage handler for vessel v1 version by HTTP GET
func V1GETStageHandler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

// V1DELETEStageHandler stage handler for vessel v1 version by HTTP DELETE
func V1DELETEStageHandler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}
