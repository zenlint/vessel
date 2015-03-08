package models

import (
	"errors"
	"fmt"
	"sync"
	"time"

	gouuid "github.com/satori/go.uuid"
)

var (
	ErrFlowNotExist = errors.New("Flow does not exist")
)

// Flow represents a complete CI solution.
type Flow struct {
	UUID string `json:"-"`
	Name string

	// Pipelines stores UUIDs of pipelines.
	Pipelines map[string]bool

	Created time.Time
}

func NewFlow(uuid, name string) *Flow {
	if len(uuid) == 0 {
		uuid = gouuid.NewV4().String()
	}
	return &Flow{
		UUID:      uuid,
		Name:      name,
		Pipelines: make(map[string]bool),
		Created:   time.Now(),
	}
}

func DeleteFlow(uuid string) error {
	return Delete(uuid, SET_TYPE_FLOW)
}

func (f *Flow) Save() error {
	return Save(f.UUID, f)
}

func (f *Flow) Retrieve() error {
	return Retrieve(f.UUID, f)
}

func (f *Flow) SetPipelines(uuids ...string) (err error) {
	var (
		pipelines = make(map[string]bool)
		pipe      *Pipeline
	)
	for _, uuid := range uuids {
		if pipelines[uuid] {
			continue
		}

		pipe = NewPipeline(uuid, "")
		if err = pipe.Retrieve(); err != nil {
			if err == ErrObjectNotExist {
				return ErrPipelineNotExist{uuid}
			} else {
				return err
			}
		}
		pipelines[uuid] = true
	}

	f.Pipelines = pipelines
	return nil
}

func ListFlows() ([]*Flow, error) {
	keys, err := LedisDB.HKeys([]byte(SET_TYPE_FLOW))
	if err != nil {
		return nil, err
	}

	flows := make([]*Flow, len(keys))
	for i := range keys {
		flows[i] = NewFlow(string(keys[i]), "")
		if err = flows[i].Retrieve(); err != nil {
			return nil, fmt.Errorf("Retrieve '%s': %v", flows[i].UUID, err)
		}
	}

	return flows, nil
}

// ___________.__                .___                 __
// \_   _____/|  |   ______  _  _|   | ____   _______/  |______    ____   ____  ____
//  |    __)  |  |  /  _ \ \/ \/ /   |/    \ /  ___/\   __\__  \  /    \_/ ___\/ __ \
//  |     \   |  |_(  <_> )     /|   |   |  \\___ \  |  |  / __ \|   |  \  \__\  ___/
//  \___  /   |____/\____/ \/\_/ |___|___|  /____  > |__| (____  /___|  /\___  >___  >
//      \/                                \/     \/            \/     \/     \/    \/

// FlowInstance represents a running instance of flow.
type FlowInstance struct {
	Flow

	// pipelines stores actual objects of pipelines.
	pipelines []*Pipeline
	// tree marks if a pipeline is done.
	treeLock sync.RWMutex
	tree     map[string]bool
}

// NewInstance creates and returns a new flow instance.
func (f *Flow) NewInstance() *FlowInstance {
	fi := &FlowInstance{
		Flow: Flow{
			UUID:    gouuid.NewV4().String(),
			Name:    f.Name,
			Created: time.Now(),
		},
		tree: make(map[string]bool),
	}

	// TODO: clone new instances of pipelines.

	return fi
}
