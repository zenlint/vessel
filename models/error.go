package models

import (
	"fmt"
)

type ErrPipelineNotExist struct {
	UUID string
}

func IsErrPipelineNotExist(err error) bool {
	_, ok := err.(ErrPipelineNotExist)
	return ok
}

func (err ErrPipelineNotExist) Error() string {
	return fmt.Sprintf("pipeline '%s' does not exist", err.UUID)
}

type ErrCircularDependencies struct {
	ObjA, ObjB string
}

func IsErrCircularDependencies(err error) bool {
	_, ok := err.(ErrCircularDependencies)
	return ok
}

func (err ErrCircularDependencies) Error() string {
	return fmt.Sprintf("circular dependencies between '%s' and '%s'", err.ObjA, err.ObjB)
}

type ErrStageNotExist struct {
	UUID string
}

func IsErrStageNotExist(err error) bool {
	_, ok := err.(ErrStageNotExist)
	return ok
}

func (err ErrStageNotExist) Error() string {
	return fmt.Sprintf("stage '%s' does not exist", err.UUID)
}

type ErrJobNotExist struct {
	UUID string
}

func IsErrJobNotExist(err error) bool {
	_, ok := err.(ErrJobNotExist)
	return ok
}

func (err ErrJobNotExist) Error() string {
	return fmt.Sprintf("job '%s' does not exist", err.UUID)
}
