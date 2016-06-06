package handler

import (
	"net/http"

	"gopkg.in/macaron.v1"
)

// PointPOSTJSON
type PointPOSTJSON struct {
}

// V1POSTPointHandler
func V1POSTPointHandler(ctx *macaron.Context, point PointPOSTJSON) (int, []byte) {
	return http.StatusOK, []byte("")
}

// V1PUTPointHandler
func V1PUTPointHandler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

// V1GETPointHandler
func V1GETPointHandler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

// V1DELETEPointHandler
func V1DELETEPointHandler(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}
