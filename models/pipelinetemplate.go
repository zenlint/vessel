package models

// PipelineSpecTemplate struct convert from for user request data
type PipelineSpecTemplate struct {
	Kind       string           `json:"kind"`
	APIVersion string           `json:"apiVersion"`
	MetaData   PipelineMetaData `json:"metadata"`
	Spec       []StageSpec      `json:"spec" sql:"-"`
}

// PipelineMetaData struct for convert from user request data
type PipelineMetaData struct {
	Name              string            `json:"name"`
	Workspace         string            `json:"workspace"`
	Project           string            `json:"project"`
	Namespace         string            `json:"namespace"`
	SelfLink          string            `json:"selfLink"`
	UID               string            `json:"uid"`
	CreateonTimestamp string            `json:"createonTimestamp"`
	DeletionTimestamp string            `json:"deletionTimestamp"`
	TimeoutDuration   int64             `json:"timeoutDuration"`
	Labels            map[string]string `json:"labels" sql:"-"`
	Annotations       map[string]string `json:"annotations" sql:"-"`
}

// StageSpec struct for convert from user request data
type StageSpec struct {
	Name                string `json:"name"`
	WorkspaceID         int64
	ProjectID           int64
	PipelineID          int64
	StageID             int64
	Replicas            int    `json:"replicas"`
	Dependence          string `json:"dependence"`
	Kind                string `json:"kind"`
	StatusCheckURL      string `json:"statusCheckLink"`
	StatusCheckInterval int64  `json:"statusCheckInterval"`
	StatusCheckCount    int    `json:"statusCheckCount"`
	Image               string `json:"image"`
	Port                int    `json:"port"`
	EnvName             string `json:"envName"`
	EnvValue            string `json:"envValue"`
}
