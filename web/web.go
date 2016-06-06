package web

import (
	"gopkg.in/macaron.v1"

	"github.com/containerops/vessel/middleware"
	_ "github.com/containerops/vessel/models"
	"github.com/containerops/vessel/router"
)

// SetVesselMacaron
func SetVesselMacaron(m *macaron.Macaron) {
	//Setting Middleware
	middleware.SetMiddlewares(m)

	//Setting Router
	router.SetRouters(m)
}
