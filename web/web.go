package web

import (
	"gopkg.in/macaron.v1"

	"github.com/containerops/vessel/middleware"
	"github.com/containerops/vessel/router"
)

func SetVesselMacaron(m *macaron.Macaron) {
	//Setting Middleware
	middleware.SetMiddlewares(m)

	//Setting Router
	router.SetRouters(m)
}
