package models

type PipelineResult struct {
	Namespace    string                `json:"-"`
	Name         string                `json:"-"`
	Status       string                `json:"status"`
	WorkspaceId  int64                 `json:"workspaceId"`
	ProjectId    int64                 `json:"projectId"`
	PipelineId   string                `json:"pipelineId"`
	Details      []*StageResult        `json:"details"`
	PipelineSpec *PipelineSpecTemplate `json:"pipelineSpec"`
}

type PipelineSpecTemplate struct {
	Kind           string           `json:"kind" binding:"Required"`
	ApiVersion     string           `json:"apiVersion" binding:"Required"`
	ApiServiceUrl  string           `json:"-"`
	ApiServiceAuth string           `json:"-"`
	MetaData       *Pipeline        `json:"metadata" binding:"Required"`
	Spec           []*Stage         `json:"spec"  binding:"Required"`
}

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