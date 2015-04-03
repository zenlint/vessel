package models

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/satori/go.uuid"

	"github.com/containerops/vessel/modules/utils"
)

// Pipeline represents a list of stages in order.
type Pipeline struct {
	UUID string `json:"-"`
	Name string

	// Stages stores UUIDs of stages.
	Stages []string
	// Requires stores prerequisites of this pipeline.
	// All of them must be done successfully in order to run this one.
	Requires map[string]bool

	Created time.Time
}

// NewPipeline creates and returns a new pipeline.
func NewPipeline(name string) *Pipeline {
	return &Pipeline{
		UUID:     uuid.NewV4().String(),
		Requires: make(map[string]bool),
		Name:     name,
		Created:  time.Now(),
		// state:  STATE_WAITING,
		// stages: make([]*Stage, 0, 3),
	}
}

func DeletePipeline(uuid string) error {
	return Delete(uuid, SET_TYPE_PIPELINE)
}

func (p *Pipeline) Save() error {
	return Save(p.UUID, p)
}

func (p *Pipeline) Retrieve() error {
	return Retrieve(p.UUID, p)
}

func (p *Pipeline) SetStages(uuids ...string) (err error) {
	var (
		marks  = make(map[string]bool)
		stage  *Stage
		stages = make([]string, 0, len(uuids))
	)
	for _, uuid := range uuids {
		if marks[uuid] {
			continue
		}

		stage = &Stage{UUID: uuid}
		if err = stage.Retrieve(); err != nil {
			if err == ErrObjectNotExist {
				return ErrStageNotExist{uuid}
			} else {
				return err
			}
		}
		marks[uuid] = true
		stages = append(stages, uuid)
	}

	p.Stages = stages
	return nil
}

func (p *Pipeline) SetPrerequisites(uuids ...string) (err error) {
	var (
		requires = make(map[string]bool)
		pipe     *Pipeline
	)
	for _, uuid := range uuids {
		if requires[uuid] || uuid == p.UUID {
			continue
		}

		pipe = &Pipeline{UUID: uuid}
		if err = pipe.Retrieve(); err != nil {
			if err == ErrObjectNotExist {
				return ErrPipelineNotExist{uuid}
			} else {
				return err
			}
		} else if pipe.Requires[p.UUID] {
			return ErrCircularDependencies{p.UUID, uuid}
		}
		requires[uuid] = true
	}

	p.Requires = requires
	return nil
}

func ListPipelines() ([]*Pipeline, error) {
	keys, err := LedisDB.HKeys([]byte(SET_TYPE_PIPELINE))
	if err != nil {
		return nil, err
	}

	pipes := make([]*Pipeline, len(keys))
	for i := range keys {
		pipes[i] = &Pipeline{UUID: string(keys[i])}
		if err = pipes[i].Retrieve(); err != nil {
			return nil, fmt.Errorf("Retrieve '%s': %v", pipes[i].UUID, err)
		}
	}

	return pipes, nil
}

// __________.__              .__  .__              .___                 __
// \______   \__|_____   ____ |  | |__| ____   ____ |   | ____   _______/  |______    ____   ____  ____
//  |     ___/  \____ \_/ __ \|  | |  |/    \_/ __ \|   |/    \ /  ___/\   __\__  \  /    \_/ ___\/ __ \
//  |    |   |  |  |_> >  ___/|  |_|  |   |  \  ___/|   |   |  \\___ \  |  |  / __ \|   |  \  \__\  ___/
//  |____|   |__|   __/ \___  >____/__|___|  /\___  >___|___|  /____  > |__| (____  /___|  /\___  >___  >
//              |__|        \/             \/     \/         \/     \/            \/     \/     \/    \/

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

// PipelineInstance represents a running pipeline.
type PipelineInstance struct {
	Pipeline

	State
	Stdout *utils.PrefixWriter

	// stages stores actual objects of stages.
	stages []*Stage
	// tree marks if a stage is done.
	treeLock sync.RWMutex
	tree     map[string]bool
}

func (p *Pipeline) NewInstance() *PipelineInstance {
	pi := &PipelineInstance{
		Pipeline: Pipeline{
			UUID:    uuid.NewV4().String(),
			Name:    p.Name,
			Created: time.Now(),
		},
		tree: make(map[string]bool),
	}

	// TODO: clone new instances of stages.
	// Clone stage UUIDs.
	pi.Stages = make([]string, len(p.Stages))
	for i := range p.Stages {
		pi.Stages[i] = p.Stages[i]
	}

	return pi
}

// AddStage adds a new stage to pipeline.
// func (p *Pipeline) AddStage(s *Stage) {
// 	for i := range p.stages {
// 		if p.stages[i].UUID() == s.UUID() {
// 			return
// 		}
// 	}
// 	p.stages = append(p.stages, s)
// }

// RemoveStage removes stage with given UUID.
// func (p *Pipeline) RemoveStage(uuid string) {
// 	for i := range p.stages {
// 		if p.stages[i].UUID() == uuid {
// 			p.stages = append(p.stages[:i], p.stages[i+1:]...)
// 			return
// 		}
// 	}
// }

// func (p *Pipeline) Run() (err error) {
// 	if p.state == STATE_RUNNING {
// 		return ErrPipelineIsBusy
// 	} else if p.state == STATE_DONE {
// 		return ErrPipelineIsDone
// 	}

// 	p.state = STATE_RUNNING
// 	defer func() {
// 		p.state = STATE_DONE
// 	}()

// 	for _, s := range p.stages {
// 		s.Stdout = utils.NewPrefixWriter("["+p.UUID+"]", p.Stdout)
// 		if err = s.Run(); err != nil {
// 			return fmt.Errorf("run stage(%s): %v", s.UUID(), err)
// 		}
// 	}

// 	return nil
// }
