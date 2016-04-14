package models

type PipelineSpecTemplate struct {
	Kind       string           `json:"kind"`
	ApiVersion string           `json:"apiVersion"`
	MetaData   PipelineMetaData `json:"metadata"`
	Spec       []StageSpec      `json:"spec"`
}

type PipelineMetaData struct {
	Name              string `json:"name"`
	Workspace         string `json:"workspace"`
	Project           string `json:"project"`
	Namespace         string `json:"namespace"`
	SelfLink          string `json:"selfLink"`
	Uid               string `json:"uid"`
	CreateonTimestamp string `json:"createonTimestamp"`
	DeletionTimestamp string `json:"deletionTimestamp"`
	TimeoutDuration   int64  `json:"timeoutDuration"`
	Labels            string `json:"labels"`
	Annotations       string `json:"annotations"`
}

type StageSpec struct {
	Name                string `json:"name"`
	Replicsa            int64  `json:"replicsa"`
	Dependence          string `json:"dependence"`
	Kind                string `json:"kind"`
	StatusCheckUrl      string `json:"statusCheckUrl"`
	StatusCheckInterval int64  `json:"statusCheckInterval"`
	StatusCheckCount    int64  `json:"statusCheckCount"`
	Image               string `json:"image"`
	Port                int64  `json:"port"`
}
