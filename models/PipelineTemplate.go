package models

type PipelineSpecTemplate struct {
	Kind       string `json:"kind"`
	ApiVersion string `json:"apiVersion"`
	MetaData   struct {
		Name              string `json:"name"`
		Workspace         string `json:"workspace"`
		Project           string `json:"project"`
		Namespace         string `json:"namespace"`
		SelfLink          string `json:"selfLink"`
		Uid               string `json:"uid"`
		CreateonTimestamp string `json:"createonTimestamp"`
		DeletionTimestamp string `json:"deletionTimestamp"`
		TimeoutDuration   string `json:"timeoutDuration"`
		Labels            string `json:"labels"`
		Annotations       string `json:"annotations"`
	} `json:"metadata"`
	Spec []struct {
		Name                string `json:"name"`
		Replicsa            int64  `json:"replicsa"`
		Dependence          string `json:"dependence"`
		Kind                string `json:"kind"`
		StatusCheckUrl      string `json:"statusCheckUrl"`
		StatusCheckInterval int64  `json:"statusCheckInterval"`
		StatusCheckCount    int64  `json:"statusCheckCount"`
		Image               string `json:"image"`
		Port                string `json:"port"`
	} `json:"spec"`
}
