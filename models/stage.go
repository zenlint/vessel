package models

// Stage struct stage data
type Stage struct {
	Namespace           string `json:"-"`
	PipelineName        string `json:"-"`
	Name                string `json:"name"  binding:"Required"`
	Replicas            uint64 `json:"replicas" binding:"Required"`
	Dependence          string `json:"dependence"`
	StatusCheckURL      string `json:"statusCheckLink"`
	StatusCheckInterval uint64 `json:"statusCheckInterval"`
	StatusCheckCount    uint64 `json:"statusCheckCount"`
	Image               string `json:"image" binding:"Required"`
	Port                uint64 `json:"port" binding:"Required"`
	EnvName             string `json:"envName"`
	EnvValue            string `json:"envValue"`
	Status              string `json:"-"`
}

// StageResult struct stage result
type StageResult struct {
	Namespace string `json:"-"`
	ID        string `json:"stageVersionID"`
	Name      string `json:"stageName"`
	Result    string `json:"runResult"`
	Detail    string `json:"detail"`
}