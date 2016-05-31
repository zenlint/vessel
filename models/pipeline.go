package models

type Pipeline struct {
	Name              string
	Namespace         string
	Stages            []string
	CreationTimestamp string
	DeletionTimestamp string
	TimeoutDuration   uint64
	Status            string
}