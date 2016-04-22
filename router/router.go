package router

import (
	"gopkg.in/macaron.v1"

	"github.com/go-macaron/binding"

	"github.com/containerops/vessel/handler"
	"github.com/containerops/vessel/models"
)

func SetRouters(m *macaron.Macaron) {
	m.Group("/v1", func() {
		m.Group("/workspace", func() {

			m.Post("/", binding.Bind(models.PipelineSpecTemplate{}), handler.V1POSTWorkspaceHandler)
			m.Put("/:workspace", binding.Bind(handler.WorkspacePUTJSON{}), handler.V1PUTWorkspaceHandler)
			m.Get("/:workspace", handler.V1GETWorkspaceHandler)
			m.Delete("/:workspace", handler.V1DELETEWorkspaceHandler)

		})

		m.Group("/project/:workspace", func() {

			m.Post("/", binding.Bind(handler.ProjectPOSTJSON{}), handler.V1POSTProjectHandler)
			m.Put("/:project", binding.Bind(handler.ProjectPUTJSON{}), handler.V1PUTProjectHandler)
			m.Get("/:project", handler.V1GETProjectHandler)
			m.Delete("/:project", handler.V1DELETEProjectHandler)

		})

		m.Group("/pipeline/:project", func() {

			m.Post("/", binding.Bind(models.PipelineSpecTemplate{}), handler.V1POSTPipelineHandler)
			m.Put("/:pipeline", handler.V1PUTPipelineHandler)
			m.Get("/:pipeline", handler.V1GETPipelineHandler)
			m.Delete("/:pipeline", binding.Bind(models.PipelineSpecTemplate{}), handler.V1DELETEPipelineHandler)

			m.Put("/:pipeline/run", handler.V1RunPipelineHandler)

			m.Get("/:pipeline/status", handler.V1StatusGETHandler)
		})

		m.Group("/point/:pipeline", func() {
			m.Post("/", binding.Bind(handler.PointPOSTJSON{}), handler.V1POSTPointHandler)
			m.Put("/:point", handler.V1PUTPointHandler)
			m.Get("/:point", handler.V1GETPointHandler)
			m.Delete("/:point", handler.V1DELETEPointHandler)
		})

		m.Group("/stage/:pipeline", func() {
			m.Post("/", binding.Bind(handler.StagePOSTJSON{}), handler.V1POSTStageHandler)
			m.Put("/:stage", handler.V1PUTStageHandler)
			m.Get("/:stage", handler.V1GETStageHandler)
			m.Delete("/:stage", handler.V1DELETEStageHandler)
		})

		m.Group("/param/:uuid", func() {
			m.Post("/", binding.Bind(handler.ParamPOSTJSON{}), handler.V1POSTParamHandler)
			m.Get("/list", handler.V1GETListParamsHandler)
			m.Put("/:param", handler.V1PUTParamHandler)
			m.Get("/:param", handler.V1GETParamHandler)
			m.Delete("/:param", handler.V1DELETEParamHandler)
		})

	})
}
