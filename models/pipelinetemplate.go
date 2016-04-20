package models

type PipelineSpecTemplate struct {
	Kind       string           `json:"kind"`
	ApiVersion string           `json:"apiVersion"`
	MetaData   PipelineMetaData `json:"metadata"`
	Spec       []StageSpec      `json:"spec"`
}

type PipelineMetaData struct {
	Name              string            `json:"name"`
	Workspace         string            `json:"workspace"`
	Project           string            `json:"project"`
	Namespace         string            `json:"namespace"`
	SelfLink          string            `json:"selfLink"`
	Uid               string            `json:"uid"`
	CreateonTimestamp string            `json:"createonTimestamp"`
	DeletionTimestamp string            `json:"deletionTimestamp"`
	TimeoutDuration   int64             `json:"timeoutDuration"`
	Labels            map[string]string `json:"labels"`
	Annotations       map[string]string `json:"annotations"`
}

/*// pipelineMetadata struct for convert from pipelineVersion.MetaData
type piplineMetadata struct {
	name            string            `json:"name, omitempty"`
	namespace       string            `json:"namespace, omitempty"`
	selfLink        string            `json:"selflink, omitempty"`
	uid             types.UID         `json:"uid, omitempty"`
	createTimestamp unversioned.Time  `json:"createTimestamp, omitempty"`
	deleteTimestamp unversioned.Time  `json:"deleteTimestamp, omitempty"`
	timeoutDuration int64             `json:"timeoutDuration, omitempty"`
	labels          map[string]string `json:"labels, omitempty"`
	annotations     map[string]string `json:"annotations, omitempty"`
}

// pipelineSpec struct for convert from pipelineVersion.Spec
type Spec struct {
	name                string `json:"name, omitempty"`
	replicas            int    `json:"replicas, omitempty"`
	dependencies        string `json:"dependencies, omitempty"`
	kind                string `json:"kind, omitempty"`
	statusCheckLink     string `json:"statusCheckLink, omitempty"`
	statusCheckInterval int64  `json:"statusCheckInterval, omitempty"`
	statusCheckCount    int64  `json:"statusCheckCount, omitempty"`
	imageName           string `json:"imagename, omitempty"`
	port                int    `json:"port, omitempty"`
}*/

type StageSpec struct {
	Name                string `json:"name"`
	WorkspaceId         int64
	ProjectId           int64
	PipelineId          int64
	StageId             int64
	Replicas            int    `json:"replicsa"`
	Dependence          string `json:"dependence"`
	Kind                string `json:"kind"`
	StatusCheckUrl      string `json:"statusCheckUrl"`
	StatusCheckInterval int64  `json:"statusCheckInterval"`
	StatusCheckCount    int    `json:"statusCheckCount"`
	Image               string `json:"image"`
	Port                int    `json:"port"`
}
