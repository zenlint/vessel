package models

import (
	"errors"
	"fmt"
	"time"

	gouuid "github.com/satori/go.uuid"

	"github.com/dockercn/vessel/modules/utils"
)

// Stage represents a build process.
type Stage struct {
	UUID string `json:"-"`
	Name string

	beforeActions []Action
	afterActions  []Action

	Created time.Time
}

// NewStage creates and returns a new stage.
func NewStage(uuid, name string) *Stage {
	if len(uuid) == 0 {
		uuid = gouuid.NewV4().String()
	}
	return &Stage{
		UUID:          uuid,
		Name:          name,
		beforeActions: make([]Action, 0, 3),
		afterActions:  make([]Action, 0, 3),
	}
}

func DeleteStage(uuid string) error {
	return Delete(uuid, SET_TYPE_STAGE)
}

func (s *Stage) Save() error {
	return Save(s.UUID, s)
}

func (s *Stage) Retrieve() error {
	return Retrieve(s.UUID, s)
}

func ListStages() ([]*Stage, error) {
	keys, err := LedisDB.HKeys([]byte(SET_TYPE_STAGE))
	if err != nil {
		return nil, err
	}

	stages := make([]*Stage, len(keys))
	for i := range keys {
		stages[i] = NewStage(string(keys[i]), "")
		if err = stages[i].Retrieve(); err != nil {
			return nil, fmt.Errorf("Retrieve '%s': %v", stages[i].UUID, err)
		}
	}

	return stages, nil
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

// func (s *Stage) SetJob(job Job) {
// 	s.job = job
// }

func doActions(acts []Action) (err error) {
	for i := range acts {
		if err = acts[i].Do(); err != nil {
			return fmt.Errorf("[%s] %v", acts[i].UUID(), err)
		}
	}
	return nil
}

// func (s *Stage) Run() (err error) {
// 	if s.job == nil {
// 		return ErrJobNotSet
// 	}

// 	if err = doActions(s.beforeActions); err != nil {
// 		return fmt.Errorf("before action: %v", err)
// 	} else if err = s.job.Run(); err != nil {
// 		return fmt.Errorf("run job: %v", err)
// 	} else if err = doActions(s.afterActions); err != nil {
// 		return fmt.Errorf("after action: %v", err)
// 	}

// 	return nil
// }

//   _________ __                        .___                 __
//  /   _____//  |______     ____   ____ |   | ____   _______/  |______    ____  ____
//  \_____  \\   __\__  \   / ___\_/ __ \|   |/    \ /  ___/\   __\__  \ _/ ___\/ __ \
//  /        \|  |  / __ \_/ /_/  >  ___/|   |   |  \\___ \  |  |  / __ \\  \__\  ___/
// /_______  /|__| (____  /\___  / \___  >___|___|  /____  > |__| (____  /\___  >___  >
//         \/           \//_____/      \/         \/     \/            \/     \/    \/

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

// StageInstance represents a running stage.
type StageInstance struct {
	Stage

	Stdout *utils.PrefixWriter

	job Job
}

func (s *Stage) NewInstance() *StageInstance {
	si := &StageInstance{
		Stage: Stage{
			UUID:    gouuid.NewV4().String(),
			Name:    s.Name,
			Created: time.Now(),
		},
	}

	return si
}
