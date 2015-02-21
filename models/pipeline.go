package models

import (
	"errors"
	"fmt"

	"github.com/satori/go.uuid"

	"github.com/dockercn/vessel/modules/utils"
)

var (
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
