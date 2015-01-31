package flow

import (
	"errors"
	"fmt"
	"sync"

	"github.com/satori/go.uuid"

	"github.com/dockercn/vessel/modules/utils"
)

var (
	ErrJobNotSet      = errors.New("Job has not been set")
	ErrPipelineIsBusy = errors.New("Pipeline is already running")
	ErrPipelineIsDone = errors.New("Pipeline has been done")
)

type State int

const (
	STATE_IDLE State = iota
	STATE_WAITING
	STATE_RUNNING
	STATE_DONE
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

// Pipeline represents a list of processes in order.
type Pipeline struct {
	uuid   string
	state  State
	stages []*Stage

	Requires []string

	Stdout *utils.PrefixWriter
}

// NewPipeline creates and returns a new pipeline.
func NewPipeline() *Pipeline {
	return &Pipeline{
		uuid:   uuid.NewV4().String(),
		state:  STATE_WAITING,
		stages: make([]*Stage, 0, 3),
	}
}

func (p *Pipeline) UUID() string {
	return p.uuid
}

// AddStage adds a new stage to pipeline.
func (p *Pipeline) AddStage(s *Stage) {
	for i := range p.stages {
		if p.stages[i].UUID() == s.UUID() {
			return
		}
	}
	p.stages = append(p.stages, s)
}

// RemoveStage removes stage with given UUID.
func (p *Pipeline) RemoveStage(uuid string) {
	for i := range p.stages {
		if p.stages[i].UUID() == uuid {
			p.stages = append(p.stages[:i], p.stages[i+1:]...)
			return
		}
	}
}

func (p *Pipeline) Run() (err error) {
	if p.state == STATE_RUNNING {
		return ErrPipelineIsBusy
	} else if p.state == STATE_DONE {
		return ErrPipelineIsDone
	}

	p.state = STATE_RUNNING
	defer func() {
		p.state = STATE_DONE
	}()

	for _, s := range p.stages {
		s.Stdout = utils.NewPrefixWriter("["+p.UUID()+"]", p.Stdout)
		if err = s.Run(); err != nil {
			return fmt.Errorf("run stage(%s): %v", s.UUID(), err)
		}
	}

	return nil
}

// Flow represents a complete CI solution.
type Flow struct {
	pipelines [][]*Pipeline

	treeLock sync.RWMutex
	tree     map[string]map[string]bool
}

func (f *Flow) Run() error {

	return nil
}
