package web

import (
	api "github.com/containerops/anchor"

	"github.com/containerops/vessel/models"
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
			Job:     stages[i].Job,
			Created: stages[i].Created,
		}
	}
	ctx.JSON(200, apiStages)
}

func setJob(stage *models.Stage, ctx *Context, job *string) bool {
	if job == nil {
		return false
	}
	if err := stage.SetJob(*job); err != nil {
		if models.IsErrJobNotExist(err) {
			ctx.Handle(422, err)
		} else {
			ctx.Handle(500, "Fail to add job to stage: '%s': %v", stage.UUID, err)
		}
		return true
	}
	return false
}

// POST /stges
func CreateStage(ctx *Context, form api.CreateStageOptions) {
	if ctx.HasError(form) {
		return
	}

	stage := models.NewStage(*form.Name)

	if setJob(stage, ctx, form.Job) {
		return
	}

	if err := stage.Save(); err != nil {
		ctx.Handle(500, "Fail to save stage '%s': %v", stage.UUID, err)
		return
	}

	ctx.AutoJSON(201, stage)
}

// GET /stages/:uuid
func GetStage(ctx *Context) {
	stage := &models.Stage{UUID: ctx.Params(":uuid")}
	if err := stage.Retrieve(); err != nil {
		if models.IsErrStageNotExist(err) {
			ctx.Handle(404, err)
		} else {
			ctx.Handle(500, "Fail to retrieve stage '%s': %v", stage.UUID, err)
		}
		return
	}

	ctx.AutoJSON(200, stage)
}

// POST /stages/:uuid
func UpdateStage(ctx *Context, form api.CreateStageOptions) {
	if ctx.HasError(form) {
		return
	}

	stage := &models.Stage{UUID: ctx.Params(":uuid")}
	if err := stage.Retrieve(); err != nil {
		if models.IsErrStageNotExist(err) {
			ctx.Handle(404, err)
		} else {
			ctx.Handle(500, "Fail to retrieve stage '%s': %v", stage.UUID, err)
		}
		return
	}
	if form.Name != nil {
		stage.Name = *form.Name
	}

	if setJob(stage, ctx, form.Job) {
		return
	}

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
			ctx.Handle(500, "Fail to delete stage '%s': %v", uuid, err)
			return
		}
	}
	ctx.WriteHeader(200)
}
