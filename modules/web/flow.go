package web

import (
	// "fmt"

	api "github.com/dockercn/anchor"

	"github.com/dockercn/vessel/models"
	"github.com/dockercn/vessel/modules/utils"
)

// GET /flows
func Flows(ctx *Context) {
	flows, err := models.ListFlows()
	if err != nil {
		ctx.Handle(500, "Fail to list flows: %v", err)
		return
	}

	apiFlows := make([]*api.Flow, len(flows))
	for i := range flows {
		apiFlows[i] = &api.Flow{
			UUID:      flows[i].UUID,
			Name:      flows[i].Name,
			Pipelines: utils.MapToStrings(flows[i].Pipelines),
			Created:   flows[i].Created,
		}
	}
	ctx.JSON(200, apiFlows)
}

func setPipelines(flow *models.Flow, ctx *Context, pipelines []string) bool {
	if err := flow.SetPipelines(pipelines...); err != nil {
		if models.IsErrPipelineNotExist(err) {
			ctx.Handle(422, err)
		} else {
			ctx.Handle(500, "Fail to add pipelines to flow '%s': %v", flow.UUID, err)
		}
		return true
	}
	return false
}

// POST /flows
func CreateFlow(ctx *Context, form api.CreateFlowOptions) {
	if ctx.HasError(form) {
		return
	}

	flow := models.NewFlow("", *form.Name)
	if setPipelines(flow, ctx, form.Pipelines) {
		return
	}

	if err := flow.Save(); err != nil {
		ctx.Handle(500, "Fail to save flow '%s': %v", flow.UUID, err)
		return
	}

	ctx.AutoJSON(201, flow)
}

// GET /flows/:uuid
func GetFlow(ctx *Context) {
	flow := models.NewFlow(ctx.Params(":uuid"), "")
	if err := flow.Retrieve(); err != nil {
		if err == models.ErrObjectNotExist {
			ctx.Handle(404, models.ErrFlowNotExist)
		} else {
			ctx.Handle(500, "Fail to retrieve flow '%s': %v", flow.UUID, err)
		}
		return
	}

	ctx.AutoJSON(200, flow)
}

// POST /flows/:uuid
func UpdateFlow(ctx *Context, form api.CreateFlowOptions) {
	if ctx.HasError(form) {
		return
	}

	flow := models.NewFlow(ctx.Params(":uuid"), "")
	if err := flow.Retrieve(); err != nil {
		if err == models.ErrObjectNotExist {
			ctx.Handle(404, models.ErrFlowNotExist)
		} else {
			ctx.Handle(500, "Fail to retrieve flow '%s': %v", flow.UUID, err)
		}
		return
	}
	flow.Name = *form.Name

	if setPipelines(flow, ctx, form.Pipelines) {
		return
	}

	if err := flow.Save(); err != nil {
		ctx.Handle(500, "Fail to save flow '%s': %v", flow.UUID, err)
		return
	}

	ctx.AutoJSON(201, flow)
}

// DELETE /flows/:uuid
func DeleteFlow(ctx *Context) {
	uuid := ctx.Params(":uuid")
	if err := models.DeleteFlow(uuid); err != nil {
		if err != models.ErrObjectNotExist {
			ctx.Handle(500, "Fail to delete flow '%s': %v", uuid, err)
			return
		}
	}
	ctx.WriteHeader(200)
}
