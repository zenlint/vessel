package models

type PipelineSpecTemplate struct {
	Kind       string            `json:"kind"`
	ApiVersion string            `json:"apiVersion"`
	MetaData   *PipelineMetaData `json:"metadata"`
	Spec       []*StageSpec     `json:"spec" sql:"-"`
}

type PipelineMetaData struct {
	Name              string    `json:"name"`
	Namespace         string    `json:"namespace"`
	CreationTimestamp string    `json:"creationTimestamp"`
	DeletionTimestamp string    `json:"deletionTimestamp"`
	TimeoutDuration   uint64    `json:"timeoutDuration"`
}

type StageSpec struct {
	Name                string `json:"name"`
	Replicas            uint   `json:"replicas"`
	Dependence          string `json:"dependence"`
	StatusCheckUrl      string `json:"statusCheckLink"`
	StatusCheckInterval uint64 `json:"statusCheckInterval"`
	StatusCheckCount    uint64 `json:"statusCheckCount"`
	Image               string `json:"image"`
	Port                uint64 `json:"port"`
	EnvName             string `json:"envName"`
	EnvValue            string `json:"envValue"`
}
