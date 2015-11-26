package router

import (
	"gopkg.in/macaron.v1"

	"github.com/macaron-contrib/binding"

	"github.com/containerops/vessel/handler"
)

func SetRouters(m *macaron.Macaron) {
	m.Group("/v1", func() {
		m.Group("/workspace", func() {

			m.Post("/", binding.Bind(handler.WorkspacePOSTJSON{}), handler.V1POSTWorkspaceHandler)
			m.Put("/:workspace", handler.V1PUTWorkspaceHandler)
			m.Get("/:workspace", handler.V1GETWorkspaceHandler)
			m.Delete("/:workspace", handler.V1DELETEWorkspaceHandler)

			m.Group("/:workspace/project", func() {

				m.Post("/", binding.Bind(handler.ProjectPOSTJSON{}), handler.V1POSTProjectHandler)
				m.Put("/:project", handler.V1PUTProjectHandler)
				m.Get("/:project", handler.V1GETProjectHandler)
				m.Delete("/:project", handler.V1DELETEProjectHandler)

				m.Group("/:project/pipeline", func() {

					m.Post("/", binding.Bind(handler.PipelinePOSTJSON{}), handler.V1POSTPipelineHandler)
					m.Put("/:pipeline", handler.V1PUTPipelineHandler)
					m.Get("/:pipeline", handler.V1GETPipelineHandler)
					m.Delete("/:pipeline", handler.V1DELETEPipelineHandler)

					m.Put("/:pipeline/run", handler.V1RunPipelineHandler)

					m.Get("/:pipeline/status", handler.V1StatusGETHandler)

					m.Group("/:pipeline/point", func() {
						m.Post("/", binding.Bind(handler.PointPOSTJSON{}), handler.V1POSTPointHandler)
						m.Put("/:point", handler.V1PUTPointHandler)
						m.Get("/:point", handler.V1GETPointHandler)
						m.Delete("/:point", handler.V1DELETEPointHandler)
					})

					m.Group("/:pipeline/stage", func() {
						m.Post("/", binding.Bind(handler.StagePOSTJSON{}), handler.V1POSTStageHandler)
						m.Put("/:point", handler.V1PUTStageHandler)
						m.Get("/:point", handler.V1GETStageHandler)
						m.Delete("/:point", handler.V1DELETEStageHandler)
					})

					m.Group("/:pipeline/param", func() {
						m.Post("/", binding.Bind(handler.ParamPOSTJSON{}), handler.V1POSTParamHandler)
						m.Put("/:point", handler.V1PUTParamHandler)
						m.Get("/:point", handler.V1GETParamHandler)
						m.Delete("/:point", handler.V1DELETEParamHandler)
					})

				})
			})
		})

	})
}
