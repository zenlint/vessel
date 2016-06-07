package models

// PipelineResult struct pipeline result
type PipelineResult struct {
	Namespace      string                `json:"-"`
	Name           string                `json:"-"`
	Status         string                `json:"status"`
	WorkspaceID    int64                 `json:"workspaceId"`
	ProjectID      int64                 `json:"projectId"`
	PipelineID     string                `json:"pipelineId"`
	PipelineDetail string                `json:"pipelineDetail"`
	Details        []interface{}         `json:"details"`
	PipelineSpec   *PipelineSpecTemplate `json:"pipelineSpec"`
}

// PipelineSpecTemplate struct template for request data
type PipelineSpecTemplate struct {
	Kind           string           `json:"kind" binding:"In(CCloud)"`
	APIVersion     string           `json:"apiVersion" binding:"In(v1)"`
	APIServiceURL  string           `json:"-"`
	APIServiceAuth string           `json:"-"`
	MetaData       *Pipeline        `json:"metadata" binding:"Required"`
	Spec           []*Stage         `json:"spec"  binding:"Required"`
}

// Pipeline struct pipeline data
type Pipeline struct {
	Namespace         string            `json:"namespace" binding:"Required"`
	Name              string            `json:"name" binding:"Required"`
	SelfLink          string            `json:"-"`
	CreationTimestamp string            `json:"creationTimestamp"`
	DeletionTimestamp string            `json:"deletionTimestamp"`
	TimeoutDuration   uint64            `json:"timeoutDuration"`
	Labels            map[string]string `json:"-"`
	Annotations       map[string]string `json:"-"`
	Stages            []string          `json:"-"`
	Status            string            `json:"-"`
}