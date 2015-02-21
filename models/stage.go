package models

import (
	"errors"
	"fmt"

	"github.com/satori/go.uuid"

	"github.com/dockercn/vessel/modules/utils"
)

var (
	ErrJobNotSet = errors.New("Job has not been set")
)

// Action represents a build action before or after stage.
type Action interface {
	UUID() string
	Do() error
}

// Job is the interface that can run as a function.
type Job interface {
	Run() error
}

// Stage represents a build process.
type Stage struct {
	uuid string
	job  Job

	beforeActions []Action
	afterActions  []Action

	Stdout *utils.PrefixWriter
}

// NewStage creates and returns a new stage.
func NewStage() *Stage {
	return &Stage{
		uuid:          uuid.NewV4().String(),
		beforeActions: make([]Action, 0, 3),
		afterActions:  make([]Action, 0, 3),
	}
}

func (s *Stage) UUID() string {
	return s.uuid
}

func addAction(acts []Action, act Action) []Action {
	for i := range acts {
		if acts[i].UUID() == act.UUID() {
			return acts
		}
	}
	return append(acts, act)
}

func (s *Stage) AddBeforeAction(act Action) {
	s.beforeActions = addAction(s.beforeActions, act)
}

func (s *Stage) AddAfterAction(act Action) {
	s.afterActions = addAction(s.afterActions, act)
}

func removeAction(acts []Action, uuid string) []Action {
	for i := range acts {
		if acts[i].UUID() == uuid {
			return append(acts[:i], acts[i+1:]...)
		}
	}
	return acts
}

func (s *Stage) RemoveBeforeAction(uuid string) {
	s.beforeActions = removeAction(s.beforeActions, uuid)
}

func (s *Stage) RemoveAfterAction(uuid string) {
	s.afterActions = removeAction(s.afterActions, uuid)
}

func (s *Stage) SetJob(job Job) {
	s.job = job
}

func doActions(acts []Action) (err error) {
	for i := range acts {
		if err = acts[i].Do(); err != nil {
			return fmt.Errorf("[%s] %v", acts[i].UUID(), err)
		}
	}
	return nil
}

func (s *Stage) Run() (err error) {
	if s.job == nil {
		return ErrJobNotSet
	}

	if err = doActions(s.beforeActions); err != nil {
		return fmt.Errorf("before action: %v", err)
	} else if err = s.job.Run(); err != nil {
		return fmt.Errorf("run job: %v", err)
	} else if err = doActions(s.afterActions); err != nil {
		return fmt.Errorf("after action: %v", err)
	}

	return nil
}
