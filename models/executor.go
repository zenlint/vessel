package models

const (
	// StateReadying  operating state readying
	StateReadying = "Readying"
	// StateRunning  operating state running
	StateRunning = "Running"
	// StateDeleted  operating state deleted
	StateDeleted = "Deleted"

	// ResultSuccess  operating result success
	ResultSuccess = "OK"
	// ResultFailed  operating result failed
	ResultFailed = "Error"
	// ResultTimeout  operating result timeout
	ResultTimeout = "Timeout"
)

// Executor struct with pre condition and info
type Executor struct {
	From []string
	Info interface{}
}

// SchedulingResult struct scheduling result
type SchedulingResult struct {
	Name   string
	Status string
	Result interface{}
}
