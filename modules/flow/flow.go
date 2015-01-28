package flow

// Trigger is the interface that contains conditional channel to trigger pipeline or stage.
type Trigger interface {
}

// Stage represents a build process.
type Stage struct {
	Name string

	beforeTrigger Trigger
	afterTrigger  Trigger
}

// Pipeline represents a list of processes in order.
type Pipeline struct {
	Name string

	stages         []*Stage
	beforeTriggers []Trigger
	afterTriggers  []Trigger
}

// AddStage adds a new stage to pipeline.
func (p *Pipeline) AddStage(s *Stage) error {
	// TODO(unknwon):
	return nil
}

// RemoveStage removes stage with given name.
func (p *Pipeline) RemoveStage(name string) error {
	// TODO(unknwon):
	return nil
}

// Flow represents a complete CI process.
type Flow struct {
	pipelines []*Pipeline
}

// AddPipeline adds a new pipeline to flow.
func (f *Flow) AddPipeline(p *Pipeline) error {
	// TODO(unknwon):
	return nil
}

// RemovePipeline removes pipeline with given name.
func (f *Flow) RemovePipeline(name string) error {
	// TODO(unknwon):
	return nil
}
