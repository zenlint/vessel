package web

import (
	api "github.com/dockercn/anchor"

	"github.com/dockercn/vessel/models"
)

// GET /stages
func Stages(ctx *Context) {
	stages, err := models.ListStages()
	if err != nil {
		ctx.Handle(500, "Fail to list stages: %v", err)
		return
	}

	apiStages := make([]*api.Stage, len(stages))
	for i := range stages {
		apiStages[i] = &api.Stage{
			UUID:    stages[i].UUID,
			Name:    stages[i].Name,
			Created: stages[i].Created,
		}
	}
	ctx.JSON(200, apiStages)
}

// POST /stges
func CreateStage(ctx *Context, form api.CreateStageOptions) {
	if ctx.HasError(form) {
		return
	}

	stage := models.NewStage("", *form.Name)
	if err := stage.Save(); err != nil {
		ctx.Handle(500, "Fail to save stage '%s': %v", stage.UUID, err)
		return
	}

	ctx.AutoJSON(201, stage)
}

// POST /stages/:uuid
func UpdateStage(ctx *Context, form api.CreateStageOptions) {
	if ctx.HasError(form) {
		return
	}

	stage := models.NewStage(ctx.Params(":uuid"), "")
	if err := stage.Retrieve(); err != nil {
		if models.IsErrStageNotExist(err) {
			ctx.Handle(404, err)
		} else {
			ctx.Handle(500, "Fail to retrieve stage '%s': %v", stage.UUID, err)
		}
		return
	}
	stage.Name = *form.Name

	if err := stage.Save(); err != nil {
		ctx.Handle(500, "Fail to save stage '%s': %v", stage.UUID, err)
		return
	}

	ctx.AutoJSON(201, stage)
}

// DELETE /stages/:uuid
func DeleteStage(ctx *Context) {
	uuid := ctx.Params(":uuid")
	if err := models.DeleteStage(uuid); err != nil {
		if err != models.ErrObjectNotExist {
			ctx.Handle(500, "Fail to delete pipeline '%s': %v", uuid, err)
			return
		}
	}
	ctx.WriteHeader(200)
}
