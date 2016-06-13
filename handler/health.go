package handler

import (
	"gopkg.in/macaron.v1"
	"net/http"
)

// V1GETHealth handler for HTTP GET
func V1GETHealth(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("0")
}