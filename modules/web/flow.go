package web

import (
	// "fmt"

	api "github.com/dockercn/anchor"

	"github.com/dockercn/vessel/models"
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
			UUID:    flows[i].UUID,
			Name:    flows[i].Name,
			Created: flows[i].Created,
		}
	}
	ctx.JSON(200, apiFlows)
}

// POST /flows
func CreateFlow(ctx *Context, form api.CreateFlowOptions) {
	if ctx.HasError(form) {
		return
	}

	flow := models.CreateFlow("", *form.Name)
	if err := flow.Save(); err != nil {
		ctx.Handle(500, "Fail to save flow '%s': %v", flow.UUID, err)
		return
	}

	ctx.AutoJSON(200, flow)
}

// GET /flows/:uuid
func GetFlow(ctx *Context) {
	flow := models.CreateFlow(ctx.Params(":uuid"), "")
	if err := flow.Retrieve(); err != nil {
		if err == models.ErrObjectNotExist {
			ctx.Handle(422, "Flow does not exist")
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

	flow := models.CreateFlow(ctx.Params(":uuid"), "")
	if err := flow.Retrieve(); err != nil {
		if err == models.ErrObjectNotExist {
			ctx.Handle(422, "Flow does not exist")
		} else {
			ctx.Handle(500, "Fail to retrieve flow '%s': %v", flow.UUID, err)
		}
		return
	}
	flow.Name = *form.Name
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
