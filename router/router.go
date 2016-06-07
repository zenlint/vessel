package router

import (
	"gopkg.in/macaron.v1"
)

// SetRouters set routers for vessel http client
func SetRouters(m *macaron.Macaron) {
	m.Group("/v1", func() {

	})
}
