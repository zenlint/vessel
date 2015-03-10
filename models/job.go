package models

import (
	"fmt"
	"time"

	"github.com/satori/go.uuid"
)

type Job struct {
	UUID    string `json:"-"`
	Name    string
	Content string `json:"-"`
	Created time.Time
}

// NewJob creates and returns a new job.
func NewJob(name string) *Job {
	return &Job{
		UUID:    uuid.NewV4().String(),
		Name:    name,
		Created: time.Now(),
	}
}

func DeleteJob(uuid string) error {
	return Delete(uuid, SET_TYPE_JOB)
}

func (j *Job) Save() error {
	return Save(j.UUID, j)
}

func (j *Job) Retrieve() error {
	return Retrieve(j.UUID, j)
}

func ListJobs() ([]*Job, error) {
	keys, err := LedisDB.HKeys([]byte(SET_TYPE_JOB))
	if err != nil {
		return nil, err
	}

	jobs := make([]*Job, len(keys))
	for i := range keys {
		jobs[i] = &Job{UUID: string(keys[i])}
		if err = jobs[i].Retrieve(); err != nil {
			return nil, fmt.Errorf("Retrieve '%s': %v", jobs[i].UUID, err)
		}
	}

	return jobs, nil
}
