package web

import (
	api "github.com/dockercn/anchor"

	"github.com/dockercn/vessel/models"
	"github.com/dockercn/vessel/modules/utils"
)

// GET /pipelines
func Pipelines(ctx *Context) {
	pipes, err := models.ListPipelines()
	if err != nil {
		ctx.Handle(500, "Fail to list pipelines: %v", err)
		return
	}

	apiPipes := make([]*api.Pipeline, len(pipes))
	for i := range pipes {
		apiPipes[i] = &api.Pipeline{
			UUID:     pipes[i].UUID,
			Name:     pipes[i].Name,
			Requires: utils.MapToStrings(pipes[i].Requires),
			Created:  pipes[i].Created,
		}
	}
	ctx.JSON(200, apiPipes)
}

func setStages(pipe *models.Pipeline, ctx *Context, stages []string) bool {
	if err := pipe.SetStages(stages...); err != nil {
		if models.IsErrStageNotExist(err) {
			ctx.Handle(422, err)
		} else {
			ctx.Handle(500, "Fail to add stages to pipeline: '%s': %v", pipe.UUID, err)
		}
		return true
	}
	return false
}

func setPrerequisites(pipe *models.Pipeline, ctx *Context, requires []string) bool {
	if err := pipe.SetPrerequisites(requires...); err != nil {
		if models.IsErrPipelineNotExist(err) ||
			models.IsErrCircularDependencies(err) {
			ctx.Handle(422, err)
		} else {
			ctx.Handle(500, "Fail to add prerequisites to pipeline '%s': %v", pipe.UUID, err)
		}
		return true
	}
	return false
}

// POST /pipelines
func CreatePipeline(ctx *Context, form api.CreatePipelineOptions) {
	if ctx.HasError(form) {
		return
	}

	pipe := models.NewPipeline(*form.Name)
	if setPrerequisites(pipe, ctx, form.Requires) {
		return
	} else if setStages(pipe, ctx, form.Stages) {
		return
	}

	if err := pipe.Save(); err != nil {
		ctx.Handle(500, "Fail to save pipeline '%s': %v", pipe.UUID, err)
		return
	}

	ctx.AutoJSON(201, pipe)
}

// GET /pipelines/:uuid
func GetPipeline(ctx *Context) {
	pipe := &models.Pipeline{UUID: ctx.Params(":uuid")}
	if err := pipe.Retrieve(); err != nil {
		if models.IsErrPipelineNotExist(err) {
			ctx.Handle(404, err)
		} else {
			ctx.Handle(500, "Fail to retrieve pipeline '%s': %v", pipe.UUID, err)
		}
		return
	}

	ctx.AutoJSON(200, pipe)
}

// POST /pipelines/:uuid
func UpdatePipeline(ctx *Context, form api.CreatePipelineOptions) {
	if ctx.HasError(form) {
		return
	}

	pipe := &models.Pipeline{UUID: ctx.Params(":uuid")}
	if err := pipe.Retrieve(); err != nil {
		if models.IsErrPipelineNotExist(err) {
			ctx.Handle(404, err)
		} else {
			ctx.Handle(500, "Fail to retrieve pipeline '%s': %v", pipe.UUID, err)
		}
		return
	}
	pipe.Name = *form.Name

	if setPrerequisites(pipe, ctx, form.Requires) {
		return
	} else if setStages(pipe, ctx, form.Stages) {
		return
	}

	if err := pipe.Save(); err != nil {
		ctx.Handle(500, "Fail to save pipeline '%s': %v", pipe.UUID, err)
		return
	}

	ctx.AutoJSON(201, pipe)
}

// DELETE /pipelines/:uuid
func DeletePipeline(ctx *Context) {
	uuid := ctx.Params(":uuid")
	if err := models.DeletePipeline(uuid); err != nil {
		if err != models.ErrObjectNotExist {
			ctx.Handle(500, "Fail to delete pipeline '%s': %v", uuid, err)
			return
		}
	}
	ctx.WriteHeader(200)
}
