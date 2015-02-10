package flow

import (
	"sync"
)

// Flow represents a complete CI solution.
type Flow struct {
	pipelines [][]*Pipeline

	treeLock sync.RWMutex
	tree     map[string]map[string]bool
}

func (f *Flow) Run() error {

	return nil
}
