package kubernetes

import (
	"time"

	"github.com/containerops/vessel/models"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/watch"
	"fmt"
	"github.com/containerops/vessel/utils"
)

func CreateRC(piplelineVersion *models.PipelineSpecTemplate, stageName string) error {
	stagespecs := piplelineVersion.Spec
	metadata := piplelineVersion.MetaData
	for _, stagespec := range stagespecs {
		if stagespec.Name == stageName {
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
			if stagespec.EnvName != "" && stagespec.EnvValue != "" {
				rc.Spec.Template.Spec.Containers[0] = api.Container{Ports: []api.ContainerPort{api.ContainerPort{
					Name:          stagespec.Name,
					ContainerPort: stagespec.Port}},
					Name:            stagespec.Name,
					Image:           stagespec.Image,
					ImagePullPolicy: "IfNotPresent",
					Env: []api.EnvVar{api.EnvVar{
						Name:  stagespec.EnvName,
						Value: stagespec.EnvValue}},
				}
			} else {
				rc.Spec.Template.Spec.Containers[0] = api.Container{Ports: []api.ContainerPort{api.ContainerPort{
					Name:          stagespec.Name,
					ContainerPort: stagespec.Port}},
					Name:            stagespec.Name,
					Image:           stagespec.Image,
					ImagePullPolicy: "IfNotPresent",
				}
			}
			/*
				if stagespec.EnvName != "" && stagespec.EnvValue != "" {
					rc.Spec.Template.Spec.Containers[0].Env
				}*/
			rc.Spec.Selector["app"] = stagespec.Name

			if _, err := models.K8sClient.ReplicationControllers(metadata.Namespace).Create(rc); err != nil {
				fmt.Println(utils.CurrentLocation(), "Create rc err : ", err)
				return err
			}
		}
	}
	return nil
}

func DeleteRC(pipelineVersion *models.PipelineSpecTemplate) error {
	return nil
}

// WatchServiceStatus return status of the operation(specified by checkOp) of the pod, OK, TIMEOUT.
func WatchRCStatus(Namespace string, labelKey string, labelValue string, timeout int64, checkOp string, ch chan string) {
	fmt.Println(utils.CurrentLocation(), "Enter WatchRCStatus")
	if checkOp != string(watch.Deleted) && checkOp != string(watch.Added) {
		fmt.Println(utils.CurrentLocation(), "Params checkOp err, checkOp: ", checkOp)
	}

	opts := api.ListOptions{LabelSelector: labels.Set{labelKey: labelValue}.AsSelector()}
	w, err := models.K8sClient.ReplicationControllers(Namespace).Watch(opts)
	if err != nil {
		ch <- Error
		fmt.Println(utils.CurrentLocation(), "Get watch RC interface err ", err)
		return
	}

	t := time.NewTimer(time.Second * time.Duration(timeout))
	select {
	case event, ok := <-w.ResultChan():
		if !ok {
			fmt.Println(utils.CurrentLocation(), "Get RC event !ok")
			ch <- Error
		} else if string(event.Type) == checkOp {
			fmt.Println(utils.CurrentLocation(), "Get RC event ok")
			ch <- OK
		}

	case <-t.C:
		fmt.Println(utils.CurrentLocation(), "WatchRCStatus timeout")
		ch <- Timeout
	}
}

// checkRC rc have no status, once the rc are found, it is with running status
func CheckRC(namespace string, rcName string) bool {
	rcs, err := models.K8sClient.ReplicationControllers(namespace).List(api.ListOptions{})
	if err != nil {
		fmt.Println(utils.CurrentLocation(), "List rcs err: ", err.Error())
	}

	for _, rc := range rcs.Items {
		if rc.Name == rcName {
			return true
		}
	}
	return false
}
