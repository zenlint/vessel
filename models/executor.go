package models

import "time"

const (
	STATE_NOT_TART 	= "NotStart"
	STATE_STARTING 	= "Working"
	STATE_SUCCESS 	= "Running"
	STATE_DELETED 	= "Deleted"

	RESULT_SUCCESS	= "OK"
	RESULT_FAILED 	= "Error"
	RESULT_TIMEOUT 	= "Timeout"
)

type Start interface {
	Start(finishChan chan Result, endTime time.Time)
	IsReady(dependenceName string) bool
}

type Result interface {
	GetResult() interface{}
	GetName() string
}