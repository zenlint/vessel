package models

const (
	STATE_NOT_TART = "NotStart"
	STATE_STARTING = "Working"
	STATE_SUCCESS = "Running"
	STATE_DELETED = "Deleted"

	RESULT_SUCCESS = "OK"
	RESULT_FAILED = "Error"
	RESULT_TIMEOUT = "Timeout"
)

type Executor struct {
	From []string
	Info interface{}
}

type SchedulingResult struct {
	Name   string
	Status string
	Result interface{}
}