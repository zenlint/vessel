package web

import (
	api "github.com/dockercn/anchor"

	"github.com/dockercn/vessel/models"
)

// GET /jobs
func Jobs(ctx *Context) {
	jobs, err := models.ListJobs()
	if err != nil {
		ctx.Handle(500, "Fail to list jobs: %v", err)
		return
	}

	apiJobs := make([]*api.Job, len(jobs))
	for i := range jobs {
		apiJobs[i] = &api.Job{
			UUID:    jobs[i].UUID,
			Name:    jobs[i].Name,
			Created: jobs[i].Created,
		}
	}
	ctx.JSON(200, apiJobs)
}

// POST /jobs
func CreateJob(ctx *Context, form api.CreateJobOptions) {
	if ctx.HasError(form) {
		return
	}

	job := models.NewJob(*form.Name)
	if err := job.Save(); err != nil {
		ctx.Handle(500, "Fail to save job '%s': %v", job.UUID, err)
		return
	}

	ctx.AutoJSON(201, job)
}

// GET /jobs/:uuid
func GetJob(ctx *Context) {
	job := &models.Job{UUID: ctx.Params(":uuid")}
	if err := job.Retrieve(); err != nil {
		if models.IsErrJobNotExist(err) {
			ctx.Handle(404, err)
		} else {
			ctx.Handle(500, "Fail to retrieve job '%s': %v", job.UUID, err)
		}
		return
	}

	ctx.AutoJSON(200, job)
}

// POST /jobs/:uuid
func UpdateJob(ctx *Context, form api.CreateJobOptions) {
	if ctx.HasError(form) {
		return
	}

	job := &models.Job{UUID: ctx.Params(":uuid")}
	if err := job.Retrieve(); err != nil {
		if models.IsErrJobNotExist(err) {
			ctx.Handle(404, err)
		} else {
			ctx.Handle(500, "Fail to retrieve job '%s': %v", job.UUID, err)
		}
		return
	}
	job.Name = *form.Name

	if err := job.Save(); err != nil {
		ctx.Handle(500, "Fail to save job '%s': %v", job.UUID, err)
		return
	}

	ctx.AutoJSON(201, job)
}

// DELETE /jobs/:uuid
func DeleteJob(ctx *Context) {
	uuid := ctx.Params(":uuid")
	if err := models.DeleteJob(uuid); err != nil {
		if err != models.ErrObjectNotExist {
			ctx.Handle(500, "Fail to delete job '%s': %v", uuid, err)
			return
		}
	}
	ctx.WriteHeader(200)
}
