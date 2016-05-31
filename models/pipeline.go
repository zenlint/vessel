package models

type Pipeline struct {
	Name              string
	Namespace         string
	Stages            []string
	CreationTimestamp string
	DeletionTimestamp string
	TimeoutDuration   int64
	Status            string
}