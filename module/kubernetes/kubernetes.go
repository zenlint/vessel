package kubernetes

import (
	"fmt"

	// "github.com/containerops/vessel/models"
	// "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/client/restclient"
	"k8s.io/kubernetes/pkg/client/unversioned"
	// "k8s.io/kubernetes/pkg/util/intstr"
)

const (
	Added   = "ADDED"
	Deleted = "DELETED"
)

const (
	Error   = "ERROR"
	Timeout = "TIMEOUT"
	OK      = "OK"
)

var CLIENT *unversioned.Client

func New(hostIp string) {
	clientConfig := restclient.Config{}
	clientConfig.Host = hostIp
	client, err := unversioned.New(&clientConfig)
	if err != nil {
		fmt.Errorf("New unversioned client err: %v!\n", err.Error())
	}
	// client.ConfigMaps(namespace).
	CLIENT = client
}

/*
// Lay here for back up,the func have been moved to pipeline.go, as StartPipelin
func CreateK8SResource(pipelineversion *models.PipelineVersion) error {
	piplineMetadata := pipelineversion.MetaData
	stagespecs := pipelineversion.StageSpecs
	for _, stagespec := range stagespecs {
		rc := &api.ReplicationController{
			ObjectMeta: api.ObjectMeta{
				Labels: map[string]string{},
			},
			Spec: api.ReplicationControllerSpec{
				Template: &api.PodTemplateSpec{
					ObjectMeta: api.ObjectMeta{
						Labels: map[string]string{},
					},
				},
				Selector: map[string]string{},
			},
		}

		service := &api.Service{
			ObjectMeta: api.ObjectMeta{
				Labels: map[string]string{},
			},
			Spec: api.ServiceSpec{
				Selector: map[string]string{},
			},
		}
		rc.Spec.Template.Spec.Containers = make([]api.Container, 1)
		service.Spec.Ports = make([]api.ServicePort, 1)

		rc.SetName(stagespec.Name)
		// rc.SetNamespace(api.NamespaceDefault)
		rc.SetNamespace(piplineMetadata.Namespace)
		rc.Labels["app"] = stagespec.Name
		rc.Spec.Replicas = stagespec.Replicas
		rc.Spec.Template.SetName(stagespec.Name)
		rc.Spec.Template.Labels["app"] = stagespec.Name
		rc.Spec.Template.Spec.Containers[0] = api.Container{Ports: []api.ContainerPort{api.ContainerPort{
			Name:          stagespec.Name,
			ContainerPort: stagespec.Port}},
			Name:  stagespec.Name,
			Image: stagespec.Image}
		rc.Spec.Selector["app"] = stagespec.Name

		service.ObjectMeta.SetName(stagespec.Name)
		// service.ObjectMeta.SetNamespace(api.NamespaceDefault)
		service.ObjectMeta.SetNamespace(piplineMetadata.Namespace)
		service.ObjectMeta.Labels["app"] = stagespec.Name
		service.Spec.Ports[0] = api.ServicePort{Port: stagespec.Port, TargetPort: intstr.FromString(stagespec.Name)}
		service.Spec.Selector["app"] = stagespec.Name

		// Conver to json string for debug

		// a, err := json.Marshal(rc)
		// if err != nil {
		// 	fmt.Println(err)
		// } else {
		// 	fmt.Println(string(a))
		// }

		// Going to support create namespace after we have namespace watch lib
		// _, err := CLIENT.Namespaces().Get(piplineMetadata.Namespace)
		// if err != nil {
		// 	namespaceObj := &api.Namespace{
		// 		ObjectMeta: api.ObjectMeta{Name: piplineMetadata.Namespace},
		// 	}
		// 	if _, err := CLIENT.Namespaces().Create(namespaceObj); err != nil {
		// 		fmt.Errorf("Create namespace err : %v\n", err)
		// 		return err
		// 	}
		// 	fmt.Println("dddddd")
		// }

		if _, err := CLIENT.ReplicationControllers(piplineMetadata.Namespace).Create(rc); err != nil {
			fmt.Errorf("Create rc err : %v\n", err)
			return err
		}

		if _, err := CLIENT.Services(piplineMetadata.Namespace).Create(service); err != nil {
			fmt.Errorf("Create service err : %v\n", err)
			return err
		}
	}
	return nil
}*/
