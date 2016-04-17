package kubernetes

import (
	// "encoding/json"
	"errors"
	"fmt"
	"time"

	// "k8s.io/kubernetes/pkg/api"
	"github.com/containerops/vessel/models"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/watch"
	// "k8s.io/kubernetes/pkg/util/intstr"
)

func CreateRC(piplelineVersion *models.PipelineVersion) error {
	// piplineMetadata := piplelineVersion.MetaData
	stagespecs := piplelineVersion.StageSpecs
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

		rc.Spec.Template.Spec.Containers = make([]api.Container, 1)
		rc.SetName(stagespec.Name)
		// rc.SetNamespace(piplineMetadata.Namespace)
		// rc.SetNamespace(api.NamespaceDefault)
		rc.SetNamespace("zenlin-namespace")
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

		if _, err := CLIENT.ReplicationControllers("zenlin-namespace").Create(rc); err != nil {
			fmt.Println("Create rc err : %v\n", err)
			return err
		}
	}
	return nil
}

// WatchServiceStatus return status of the operation(specified by checkOp) of the pod, OK, TIMEOUT.
func WatchRCStatus(Namespace string, labelKey string, labelValue string, timeout int, checkOp string) (string, error) {
	if checkOp != string(watch.Deleted) && checkOp != string(watch.Added) {
		fmt.Errorf("Params checkOp err, checkOp: %v", checkOp)
	}

	//opts := api.ListOptions{FieldSelector: fields.Set{"kind": "pod"}.AsSelector()}
	opts := api.ListOptions{LabelSelector: labels.Set{labelKey: labelValue}.AsSelector()}

	w, err := CLIENT.ReplicationControllers(Namespace).Watch(opts)
	if err != nil {
		return "", err
		fmt.Errorf("Get watch interface err")
	}

	t := time.NewTimer(time.Second * time.Duration(timeout))

	for {
		select {
		case event, ok := <-w.ResultChan():
			//fmt.Println(event.Type)
			if !ok {
				fmt.Errorf("Watch err\n")
				return "", errors.New("error occours from watch chanle")
			}
			//fmt.Println(event.Type)
			if string(event.Type) == checkOp {
				return "OK", nil
			}

		case <-t.C:
			return "TIMEOUT", nil
		}
	}
}

// checkRC rc have no status, once the rc are found, it is with running status
func CheckRC(namespace string, rcName string) bool {
	rcs, err := CLIENT.ReplicationControllers(namespace).List(api.ListOptions{})
	if err != nil {
		fmt.Errorf("List rcs err: %v\n", err.Error())
	}

	for _, rc := range rcs.Items {
		if rc.Name == rcName {
			return true
		}
	}
	return false
}
