package models

import (
	"sync"

	"github.com/satori/go.uuid"
)

// Flow represents a complete CI solution.
type Flow struct {
	UUID      string
	Name      string
	pipelines [][]*Pipeline

	// tree stores dependency relations of pipelines.
	treeLock sync.RWMutex
	tree     map[string]map[string]bool
}

func CreateFlow() *Flow {
	return &Flow{
		UUID: uuid.NewV4().String(),
		tree: make(map[string]map[string]bool),
	}
}

func (f *Flow) Run() error {

	return nil
}
