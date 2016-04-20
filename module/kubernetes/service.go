package kubernetes

import (
	"fmt"
	"time"

	"github.com/containerops/vessel/models"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/util/intstr"
	"k8s.io/kubernetes/pkg/watch"
)

func CreateService(pipelineVersion *models.PipelineSpecTemplate) error {
	stagespecs := pipelineVersion.Spec
	for _, stagespec := range stagespecs {
		service := &api.Service{
			ObjectMeta: api.ObjectMeta{
				Labels: map[string]string{},
			},
			Spec: api.ServiceSpec{
				Selector: map[string]string{},
			},
		}

		service.Spec.Ports = make([]api.ServicePort, 1)
		service.ObjectMeta.SetName(stagespec.Name)
		service.ObjectMeta.SetNamespace("zenlin-namespace")
		service.ObjectMeta.Labels["app"] = stagespec.Name
		service.Spec.Ports[0] = api.ServicePort{Port: stagespec.Port, TargetPort: intstr.FromString(stagespec.Name)}
		service.Spec.Selector["app"] = stagespec.Name

		if _, err := CLIENT.Services("zenlin-namespace").Create(service); err != nil {
			fmt.Println("Create service err : %v\n", err)
			return err
		}
	}
	return nil
}

func DeleteService(pipelineVersion *models.PipelineSpecTemplate) error {
	return nil
}

// WatchServiceStatus return status of the operation(specified by checkOp) of the pod, OK, TIMEOUT.
func WatchServiceStatus(Namespace string, labelKey string, labelValue string, timeout int64, checkOp string, ch chan string) {
	if checkOp != string(watch.Deleted) && checkOp != string(watch.Added) {
		fmt.Errorf("Params checkOp err, checkOp: %v", checkOp)
	}

	opts := api.ListOptions{LabelSelector: labels.Set{labelKey: labelValue}.AsSelector()}

	w, err := CLIENT.Services(Namespace).Watch(opts)
	if err != nil {
		ch <- Error
	}

	t := time.NewTimer(time.Second * time.Duration(timeout))

	select {
	case event, ok := <-w.ResultChan():
		if !ok {
			ch <- Error
		} else if string(event.Type) == checkOp {
			ch <- OK
		}

	case <-t.C:
		ch <- Timeout
	}
}

// CheckService service have no status, once the service are found, it is with running status
func CheckService(namespace string, serviceName string) bool {

	services, err := CLIENT.Services(namespace).List(api.ListOptions{})
	if err != nil {
		fmt.Errorf("List services err: %v\n", err.Error())
	}

	for _, s := range services.Items {
		if s.Name == serviceName {
			return true
		}
	}
	return false
}
