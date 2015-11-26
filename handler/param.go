package handler

import (
	"net/http"

	"gopkg.in/macaron.v1"
)

type ParamPOSTJSON struct {
}

func V1POSTParamHandler(ctx *macaron.Context, param ParamPOSTJSON) (int, []byte) {
	return http.StatusOK, []byte("")
}

func V1PUTParamHandler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

func V1GETParamHandler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

func V1DELETEParamHandler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}
