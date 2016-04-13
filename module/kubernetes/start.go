package kubernetes

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"k8s.io/kubernetes/pkg/api/v1"
	"k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/types"
	"k8s.io/kubernetes/pkg/watch"
)

// pipelineMetadata struct for convert from pipelineVersion.MetaData
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
type piplineSpec struct {
	name                string `json:"name, omitempty"`
	replicas            string `json:"replicas, omitempty"`
	dependencies        string `json:"dependencies, omitempty"`
	kind                string `json:"kind, omitempty"`
	statusCheckLink     string `json:"statusCheckLink, omitempty"`
	statusCheckInterval int64  `json:"statusCheckInterval, omitempty"`
	statusCheckCount    int64  `json:"statusCheckCount, omitempty"`
	imageName           string `json:"imagename, omitempty"`
	port                string `json:"port, omitempty"`
}

type PipelineVersion struct {
	Id            int64    `json:"id"`
	WorkspaceId   int64    `json:"workspaceId"`
	ProjectId     int64    `json:"projectId"`
	PipelineId    int64    `json:"pipelineId"`
	Namespace     string   `json:"namespace"`
	SelfLink      string   `json:"selfLink" gorm:"type:varchar(255)"`
	Created       int64    `json:"created"`
	Updated       int64    `json:"updated"`
	Labels        string   `json:"labels"`
	Annotations   string   `json:"annotations"`
	Detail        string   `json:"detail" gorm:"type:text"`
	StageVersions []string `json:"stageVersions"`
	Log           string   `json:"log" gorm:"type:text"`
	Status        int64    `json:"state"` // 0 not start    1 working    2 success     3 failed
	MetaData      string   `json:"metadata"`
	Spec          string   `json:"spec"`
}

// v1.ReplicationController.ObjectMeta
// v1.ReplicationController.Spec

/*
api.ReplicationController{
				ObjectMeta: api.ObjectMeta{
					Name: "foo",
					Labels: map[string]string{
						"foo":  "bar",
						"name": "baz",
					},
				},
*/
func StartK8S(pv *PipelineVersion) error {
	rc := &v1.ReplicationController{}
	service := &v1.Service{}

	var pvm piplineMetadata
	var pvs piplineSpec
	err := split(pv, &pvm, &pvs)
	if err != nil {
		return err
	}

	if err := convert(pvm, pvs, &rc, &service); err != nil {
		return err
	}

	rcRes, err := CLIENT.ReplicationControllers(namespace).Create(rc)
	if err != nil {
		fmt.Errorf("Create rc err : %v\n", err)
		return err
	}

	serviceRes, err := CLIENT.Services(namespace).Create(service)
	if err != nil {
		fmt.Errorf("Create service err : %v\n", err)
		return err
	}
	// writeBack(rcRes, serviceRes, &pvm, &pvs)
}

func split(pv *PipelineVersion, pvm *piplineMetadata, pvs *piplineSpec) error {
	err := json.Unmarshal(pv.ObjectMeta, pvm)
	if err != nil {
		fmt.Errorf("Unmarshal PipelineVersion.ObjectMeta err : %v\n", err)
		return "", err
	}

	err = json.Unmarshal(pv.Spec, pvs)
	if err != nil {
		fmt.Errorf("Unmarshal PipelineVersion.Spec err : %v\n", err)
		return "", err
	}
	return nil
}

func convert(piplineMetadata *piplineMetadata, piplineSpec *piplineSpec,
	rc *v1.ReplicationController, service *v1.Service) (string, error) {
	rcRes.Name = piplineSpec.name
	rcRes.Namespace = plMetadata.namespace
	//Use map["rc"] = Spec.name for temprory
	rcRes.Labels["rc"] = piplineSpec.name
	rcRes.Spec.Replicas = piplineSpec.replicas
	rcRes.Spec.Template.Name = piplineSpec.name
	rcRes.Spec.Template.Labels["pod"] = piplineSpec.name
	rcRes.Spec.Template.Namespace = piplineMetadata.namespace
	rcRes.Spec.Template.Spec.Containers[0].Name = piplineSpec.name
	rcRes.Spec.Template.Spec.Containers[0].Image = piplineSpec.imageName
	rcRes.Spec.Template.Spec.Containers[0].Ports[0].Name = piplineSpec.name
	rcRes.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort = piplineSpec.port
	rcRes.Spec.Selector["app"] = piplineSpec.name

	serviceRes.ObjectMeta.Name = piplineSpec.name
	serviceRes.ObjectMeta.Namespace = plMetadata.namespace
	serviceRes.ObjectMeta.Labels["service"] = piplineSpec.name
	serviceRes.Spec.Ports[0].Name = piplineSpec.name
	serviceRes.Spec.Ports[0].TargetPort = piplineSpec.port
	serviceRes.Spec.Selector["app"] = piplineSpec.name

	return piplineMetadata.namespace, nil
}
