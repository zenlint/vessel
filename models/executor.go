package models

const (
	StateReady 	= "Ready"
	StateStarting 	= "Working"
	StateSuccess 	= "OK"
	StateFailed 	= "Error"
	StateTimeout 	= "Timeout"
	StateDeleting 	= "Deleting"
	StateDeleted 	= "Deleted"
)

type ExecutorRes struct {
	Name   string
	Err    error
	Result string
}