package handler

import (
	"net/http"

	"gopkg.in/macaron.v1"
)

// ParamPOSTJSON struct for json data by HTTP POST
type ParamPOSTJSON struct {
}

// V1POSTParamHandler param handler for vessel v1 version by HTTP POST
func V1POSTParamHandler(ctx *macaron.Context, param ParamPOSTJSON) (int, []byte) {
	return http.StatusOK, []byte("")
}

// V1GETListParamsHandler param list handler for vessel v1 version by HTTP GET
func V1GETListParamsHandler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

// V1PUTParamHandler param handler for vessel v1 version by HTTP PUT
func V1PUTParamHandler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

// V1GETParamHandler param handler for vessel v1 version by HTTP GET
func V1GETParamHandler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

// V1DELETEParamHandler param handler for vessel v1 version by HTTP DELETE
func V1DELETEParamHandler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}
