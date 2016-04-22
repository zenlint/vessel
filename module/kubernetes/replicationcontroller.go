package kubernetes

import (
	"fmt"
	"time"

	"github.com/containerops/vessel/models"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/watch"
)

func CreateRC(piplelineVersion *models.PipelineSpecTemplate) error {
	stagespecs := piplelineVersion.Spec
	metadata := piplelineVersion.MetaData
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
		rc.SetNamespace(metadata.Namespace)
		rc.Labels["app"] = stagespec.Name
		rc.Spec.Replicas = stagespec.Replicas
		rc.Spec.Template.SetName(stagespec.Name)
		rc.Spec.Template.Labels["app"] = stagespec.Name
		rc.Spec.Template.Spec.Containers[0] = api.Container{Ports: []api.ContainerPort{api.ContainerPort{
			Name:          stagespec.Name,
			ContainerPort: stagespec.Port}},
			Name:            stagespec.Name,
			Image:           stagespec.Image,
			ImagePullPolicy: "IfNotPresent",
		}
		rc.Spec.Selector["app"] = stagespec.Name

		if _, err := CLIENT.ReplicationControllers(metadata.Namespace).Create(rc); err != nil {
			fmt.Println("Create rc err : %v\n", err)
			return err
		}
	}
	return nil
}

func DeleteRC(pipelineVersion *models.PipelineSpecTemplate) error {
	return nil
}

// WatchServiceStatus return status of the operation(specified by checkOp) of the pod, OK, TIMEOUT.
func WatchRCStatus(Namespace string, labelKey string, labelValue string, timeout int64, checkOp string, ch chan string) {
	fmt.Println("Enter WatchRCStatus")
	if checkOp != string(watch.Deleted) && checkOp != string(watch.Added) {
		fmt.Printf("Params checkOp err, checkOp: %v", checkOp)
	}

	opts := api.ListOptions{LabelSelector: labels.Set{labelKey: labelValue}.AsSelector()}
	w, err := CLIENT.ReplicationControllers(Namespace).Watch(opts)
	if err != nil {
		ch <- Error
		fmt.Printf("Get watch RC interface err %v\n", err)
		return
	}

	t := time.NewTimer(time.Second * time.Duration(timeout))
	select {
	case event, ok := <-w.ResultChan():
		if !ok {
			fmt.Println("Get RC event !ok")
			ch <- Error
		} else if string(event.Type) == checkOp {
			fmt.Println("Get RC event ok")
			ch <- OK
		}

	case <-t.C:
		fmt.Println("WatchRCStatus timeout")
		ch <- Timeout
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
