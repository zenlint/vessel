package kubernetes

import (
	// "encoding/json"
	"fmt"

	"github.com/containerops/vessel/models"
	"k8s.io/kubernetes/pkg/api"
	// "k8s.io/kubernetes/pkg/api/unversioned"
	// "k8s.io/kubernetes/pkg/types"
	"k8s.io/kubernetes/pkg/util/intstr"
)

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
type piplineSpec struct {
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

/*
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
*/
// unversioned.ReplicationController.ObjectMeta
// unversioned.ReplicationController.Spec

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
/*
func StartK8SResource(pv *PipelineVersion) error {
	rc := &api.ReplicationController{}
	service := &api.Service{}

	var pvm piplineMetadata
	var pvs piplineSpec
	err := split(pv, &pvm, &pvs)
	if err != nil {
		return err
	}

	namespace := convert(pvm, pvs, rc, service)

	if _, err = CLIENT.ReplicationControllers(namespace).Create(rc); err != nil {
		fmt.Errorf("Create rc err : %v\n", err)
		return err
	}

	if _, err := CLIENT.Services(namespace).Create(service); err != nil {
		fmt.Errorf("Create service err : %v\n", err)
		return err
	}
	// writeBack(rcRes, serviceRes, &pvm, &pvs)
	return nil
}
*/

/*
rcRes = &api.ReplicationController{
		ObjectMeta: api.ObjectMeta{
			Name:      piplineSpec.Name,
			Namespace: piplineMetadata.Namespace,
			Labels: map[string]string{
				"app": piplineSpec.Name,
			},
		},
		Spec: api.ReplicationControllerSpec{
			Replicas: piplineSpec.Replicas,
			Template: &api.PodTemplateSpec{
				ObjectMeta: api.ObjectMeta{
					Name: piplineSpec.Name,
					Labels: map[string]string{
						"app": piplineSpec.Name,
					},
					Namespace: piplineMetadata.Namespace,
				},
				Spec: api.PodSpec{
					Containers: []api.Container{
						api.Container{
							Name:  piplineSpec.Name,
							Image: piplineSpec.ImageName,
							Ports: []api.ContainerPort{
								api.ContainerPort{
									Name:          piplineSpec.Name,
									ContainerPort: piplineSpec.Port,
								},
							},
						},
					},
				},
			},
			Selector: map[string]string{
				"app": piplineSpec.Name,
			},
		},
	}
*/

/*
	serviceRes = &api.Service{
		ObjectMeta: api.ObjectMeta{
			Name:      piplineSpec.Name,
			Namespace: piplineMetadata.Namespace,
			Labels: map[string]string{
				"app": piplineSpec.Name,
			},
		},
		Spec: api.ServiceSpec{
			Ports: []api.ServicePort{
				api.ServicePort{
					Port:       piplineSpec.Port,
					TargetPort: intstr.FromString(piplineSpec.Name),
				},
			},
			Selector: map[string]string{
				"app": piplineSpec.Name,
			},
		},
	}
*/

/*
namespace := &api.Namespace{
		ObjectMeta: api.ObjectMeta{Name: "foo"},
	}
	c := &simple.Client{
		Request: simple.Request{
			Method: "POST",
			Path:   testapi.Default.ResourcePath("namespaces", "", ""),
			Body:   namespace,
		},
		Response: simple.Response{StatusCode: 200, Body: namespace},
	}

	// from the source ns, provision a new global namespace "foo"
	response, err := c.Setup(t).Namespaces().Create(namespace)
*/
func StartK8SResource(pipelineversion *models.PipelineVersion) error {

	rc := &api.ReplicationController{
		ObjectMeta: api.ObjectMeta{
			Labels: map[string]string{},
		},
		Spec: api.ReplicationControllerSpec{
			Template: &api.PodTemplateSpec{
				ObjectMeta: api.ObjectMeta{
					Labels: map[string]string{},
				},
				Spec: api.PodSpec{
					Containers: []api.Container{
						api.Container{
							Ports: []api.ContainerPort{},
						},
					},
				},
			},
		},
	}

	service := &api.Service{
		ObjectMeta: api.ObjectMeta{
			Labels: map[string]string{},
		},
		Spec: api.ServiceSpec{
			Ports:    []api.ServicePort{},
			Selector: map[string]string{},
		},
	}

	piplineMetadata := pipelineversion.MetaData
	stagespecs := pipelineversion.StageSpecs

	for _, stagespec := range stagespecs {
		// rc := api.ReplicationController{}
		// service := api.Service{}
		container := api.Container{}

		rc.SetName(stagespec.Name)
		rc.SetNamespace(piplineMetadata.Namespace)
		rc.Labels["app"] = stagespec.Name
		rc.Spec.Replicas = stagespec.Replicas
		rc.Spec.Template.SetName(stagespec.Name)
		rc.Spec.Template.Labels["app"] = stagespec.Name
		container.Ports = append(container.Ports, api.ContainerPort{Name: stagespec.Name, ContainerPort: stagespec.Port})
		container.Name = stagespec.Name
		container.Image = stagespec.Image
		rc.Spec.Template.Spec.Containers = append(rc.Spec.Template.Spec.Containers, container)
		rc.Spec.Selector["app"] = stagespec.Name

		service.ObjectMeta.SetName(stagespec.Name)
		service.ObjectMeta.SetNamespace(piplineMetadata.Namespace)
		service.ObjectMeta.Labels["app"] = stagespec.Name
		service.Spec.Ports = append(service.Spec.Ports, api.ServicePort{Port: stagespec.Port, TargetPort: intstr.FromString(stagespec.Name)})
		service.Spec.Selector["app"] = stagespec.Name

		namespace := &api.Namespace{
			ObjectMeta: api.ObjectMeta{Name: piplineMetadata.Namespace},
		}
		_, err := CLIENT.Namespaces().Create(namespace)
		if err != nil {
			fmt.Errorf("Create namespace err : %v\n", err)
		}

		if _, err = CLIENT.ReplicationControllers(piplineMetadata.Namespace).Create(rc); err != nil {
			fmt.Errorf("Create rc err : %v\n", err)
			return err
		}
		// CLIENT.ReplicationControllers(namespace)

		if _, err := CLIENT.Services(piplineMetadata.Namespace).Create(service); err != nil {
			fmt.Errorf("Create service err : %v\n", err)
			return err
		}
		// namespaces = append(namespaces, piplineMetadata.Namespace)
	}

	// namespaces := convert(pm, pss, rc, service)
	/*
		if _, err = CLIENT.ReplicationControllers(namespace).Create(rc); err != nil {
			fmt.Errorf("Create rc err : %v\n", err)
			return err
		}

		if _, err := CLIENT.Services(namespace).Create(service); err != nil {
			fmt.Errorf("Create service err : %v\n", err)
			return err
		}*/
	// writeBack(rcRes, serviceRes, &pvm, &pvs)
	return nil
}

/*func convert(piplineMetadata models.PipelineMetaData, stagespecs []models.StageSpec,
rcRes *[]api.ReplicationController, serviceRes *[]api.Service) []string {
namespaces := make([]string)

for i, stagespec := range stagespecs {
	rc := api.ReplicationController{}
	service := api.Service{}
	container := api.Container{}

	rc.SetName(stagespec.Name)
	rc.SetNamespace(piplineMetadata.Namespace)
	rc.Labels["app"] = stagespec.Name
	rc.Spec.Replicas = stagespec.Replicsas
	rc.Spec.Template.SetName(stagespec.Name)
	rc.Spec.Template.Labels["app"] = stagespec.Name
	container.Ports = append(container.Ports, api.ContainerPort{Name: stagespec.Name, ContainerPort: stagespec.Port})
	container.Name = stagespec.Name
	container.Image = stagespec.Image
	rc.Spec.Template.Spec.Containers[i] = append(rc.Spec.Template.Spec.Containers[i], container)
	rc.Spec.Selector["app"] = stagespec.Name

	service.ObjectMeta.SetName(stagespec.Name)
	service.ObjectMeta.SetNamespace(piplineMetadata.Namespace)
	service.ObjectMeta.Labels["app"] = stagespec.Name
	service.Spec.Ports = append(service.Spec.Ports, api.ServicePort{Port: stagespec.Port, TargetPort: intstr.FromString(stagespec.Name)})
	service.Spec.Selector["app"] = stagespec.Name

	*rcRes = append(*rcRes, rc)
	*serviceRes = append(*serviceRes, service)
	namespaces = append(namespaces, piplineMetadata.Namespace)
}
*/
/*rcRes.Name = piplineSpec.name

rcRes.Namespace = piplineMetadata.namespace
//Use map["rc"] = Spec.name for temprory`
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
serviceRes.ObjectMeta.Namespace = piplineMetadata.namespace
serviceRes.ObjectMeta.Labels["service"] = piplineSpec.name
serviceRes.Spec.Ports[0].Port = piplineSpec.port
serviceRes.Spec.Ports[0].TargetPort = intstr.FromString(piplineSpec.name)
serviceRes.Spec.Selector["app"] = piplineSpec.name*/

/*
func split(pv *PipelineVersion, pvm *piplineMetadata, pvs *piplineSpec) error {
	err := json.Unmarshal([]byte(pv.MetaData), pvm)
	if err != nil {
		fmt.Errorf("Unmarshal PipelineVersion.ObjectMeta err : %v\n", err)
		return err
	}

	err = json.Unmarshal([]byte(pv.Spec), pvs)
	if err != nil {
		fmt.Errorf("Unmarshal PipelineVersion.Spec err : %v\n", err)
		return err
	}
	return nil
}
*/

/*func split(pipeline PiplelineInterface, pvm *piplineMetadata, pvs *piplineSpec) error {
	err := json.Unmarshal([]byte(pipeline.GetMetadata()), pvm)
	if err != nil {
		fmt.Errorf("Unmarshal pipeline.GetMetadata() err : %v\n", err)
		return err
	}

	err = json.Unmarshal([]byte(pipeline.GetSpec()), pvs)
	if err != nil {
		fmt.Errorf("Unmarshal pipeline.GetSpec() err : %v\n", err)
		return err
	}
	return nil
}*/

/*
rcRes = &api.ReplicationController{
		ObjectMeta: api.ObjectMeta{
			Name:      piplineSpec.Name,
			Namespace: piplineMetadata.Namespace,
			Labels: map[string]string{
				"app": piplineSpec.Name,
			},
		},
		Spec: api.ReplicationControllerSpec{
			Replicas: piplineSpec.Replicas,
			Template: &api.PodTemplateSpec{
				ObjectMeta: api.ObjectMeta{
					Name: piplineSpec.Name,
					Labels: map[string]string{
						"app": piplineSpec.Name,
					},
					Namespace: piplineMetadata.Namespace,
				},
				Spec: api.PodSpec{
					Containers: []api.Container{
						api.Container{
							Name:  piplineSpec.Name,
							Image: piplineSpec.ImageName,
							Ports: []api.ContainerPort{
								api.ContainerPort{
									Name:          piplineSpec.Name,
									ContainerPort: piplineSpec.Port,
								},
							},
						},
					},
				},
			},
			Selector: map[string]string{
				"app": piplineSpec.Name,
			},
		},
	}
*/
