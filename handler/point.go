package handler

import (
	"net/http"

	"gopkg.in/macaron.v1"
)

// PointPOSTJSON struct for json data by HTTP POST
type PointPOSTJSON struct {
}

// V1POSTPointHandler point handler for vessel v1 version by HTTP POST
func V1POSTPointHandler(ctx *macaron.Context, point PointPOSTJSON) (int, []byte) {
	return http.StatusOK, []byte("")
}

// V1PUTPointHandler point handler for vessel v1 version by HTTP PUT
func V1PUTPointHandler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

// V1GETPointHandler point handler for vessel v1 version by HTTP GET
func V1GETPointHandler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

// V1DELETEPointHandler point handler for vessel v1 version by HTTP DELETE
func V1DELETEPointHandler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}
