package models

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/Unknwon/com"
	"github.com/satori/go.uuid"

	"github.com/dockercn/vessel/modules/utils"
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

func (j *Job) Save() (err error) {
	// Save content to data directory.
	if len(j.Content) > 0 {
		fpath := utils.UUIDToFilePath(j.UUID)
		if err = os.MkdirAll(path.Dir(fpath), os.ModePerm); err != nil {
			return fmt.Errorf("create directory: %v", err)
		} else if err = ioutil.WriteFile(fpath, []byte(j.Content), os.ModePerm); err != nil {
			return fmt.Errorf("write content: %v", err)
		}
	}
	return Save(j.UUID, j)
}

func (j *Job) Retrieve() error {
	// Read content to data directory.
	fpath := utils.UUIDToFilePath(j.UUID)
	if com.IsFile(fpath) {
		data, err := ioutil.ReadFile(fpath)
		if err != nil {
			return fmt.Errorf("read content: %v", err)
		}
		j.Content = string(data)
	}
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
