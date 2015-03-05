package models

import (
	"fmt"
	"sync"
	"time"

	gouuid "github.com/satori/go.uuid"
)

// Flow represents a complete CI solution.
type Flow struct {
	UUID string `json:"-"`
	Name string

	// pipelines stores execution dependency relations of pipelines.
	pipelines [][]*Pipeline
	Pipelines [][]string // Store UUIDs when save.

	// Tree marks whether prerequisites are satisfied for a pipeline to run.
	treeLock sync.RWMutex
	Tree     map[string]map[string]bool

	Created time.Time
}

func CreateFlow(uuid, name string) *Flow {
	if len(uuid) == 0 {
		uuid = gouuid.NewV4().String()
	}
	return &Flow{
		UUID:    uuid,
		Name:    name,
		Tree:    make(map[string]map[string]bool),
		Created: time.Now(),
	}
}

func DeleteFlow(uuid string) (err error) {
	if _, err = LedisDB.HDel([]byte(SET_TYPE_FLOW), []byte(uuid)); err != nil {
		return err
	}
	_, err = LedisDB.Del([]byte(uuid))
	return err
}

// RefreshPipelineUUIDs iterates pipelines and only store UUIDs to map their relations.
func (f *Flow) RefreshPipelineUUIDs() {
	f.Pipelines = make([][]string, len(f.pipelines))
	for i := range f.pipelines {
		f.Pipelines[i] = make([]string, len(f.pipelines[i]))
		for j := range f.pipelines[i] {
			f.Pipelines[i][j] = f.pipelines[i][j].UUID
		}
	}
}

func (f *Flow) Save() error {
	f.RefreshPipelineUUIDs()
	return Save(f.UUID, f)
}

func (f *Flow) Retrieve() error {
	return Retrieve(f.UUID, f)
}

func (f *Flow) Run() error {

	return nil
}

func ListFlows() ([]*Flow, error) {
	keys, err := LedisDB.HKeys([]byte(SET_TYPE_FLOW))
	if err != nil {
		return nil, err
	}

	flows := make([]*Flow, len(keys))
	for i := range keys {
		flows[i] = CreateFlow(string(keys[i]), "")
		if err = flows[i].Retrieve(); err != nil {
			return nil, fmt.Errorf("Retrieve '%s': %v", flows[i].UUID, err)
		}
	}

	return flows, nil
}
