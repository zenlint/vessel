package handler

import (
	"net/http"

	"gopkg.in/macaron.v1"
)

// ParamPOSTJSON
type ParamPOSTJSON struct {
}

// V1POSTParamHandler
func V1POSTParamHandler(ctx *macaron.Context, param ParamPOSTJSON) (int, []byte) {
	return http.StatusOK, []byte("")
}

// V1GETListParamsHandler
func V1GETListParamsHandler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

// V1PUTParamHandler
func V1PUTParamHandler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

// V1GETParamHandler
func V1GETParamHandler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

// V1DELETEParamHandler
func V1DELETEParamHandler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}
