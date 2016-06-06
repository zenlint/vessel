package router

import (
	"gopkg.in/macaron.v1"

	// "github.com/go-macaron/binding"
/*
	"github.com/containerops/vessel/handler"
	"github.com/containerops/vessel/models"*/
)

// SetRouters set routers for vessel http client
func SetRouters(m *macaron.Macaron) {
	m.Group("/v1", func() {

	})
}
