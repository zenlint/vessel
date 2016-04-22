package kubernetes

import (
	"fmt"
	"time"

	"github.com/containerops/vessel/models"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/watch"
)

func CreateNamespace(pipelineVersion *models.PipelineSpecTemplate) error {
	piplineMetadata := pipelineVersion.MetaData
	namespaceObj := &api.Namespace{
		ObjectMeta: api.ObjectMeta{
			Name:   piplineMetadata.Namespace,
			Labels: map[string]string{},
		},
	}
	namespaceObj.SetLabels(map[string]string{"app": piplineMetadata.Name})

	if _, err := CLIENT.Namespaces().Create(namespaceObj); err != nil {
		fmt.Errorf("Create namespace err : %v\n", err)
		return err
	}
	return nil
}

// WatchPodStatus return status of the operation(specified by checkOp) of the pod, OK, TIMEOUT.
func WatchNamespaceStatus(labelKey string, labelValue string, timeout int64, checkOp string, ch chan string) {
	if checkOp != string(watch.Deleted) && checkOp != string(watch.Added) {
		fmt.Errorf("Params checkOp err, checkOp: %v", checkOp)
	}

	opts := api.ListOptions{LabelSelector: labels.Set{labelKey: labelValue}.AsSelector()}
	w, err := CLIENT.Namespaces().Watch(opts)
	if err != nil {
		ch <- Error
		return
	}

	t := time.NewTimer(time.Second * time.Duration(timeout))
	select {
	case event, ok := <-w.ResultChan():
		if !ok {
			ch <- Error
		} else if string(event.Type) == checkOp {
			// Pod have phase, so we have to wait for the phase change to the right status when added
			fmt.Println(event.Object.(*api.Namespace).Status.Phase)

			if (checkOp == string(watch.Deleted)) || ((checkOp != string(watch.Deleted)) &&
				(event.Object.(*api.Namespace).Status.Phase == "Active")) {
				ch <- OK
			}
		}

	case <-t.C:
		fmt.Println("Watch namespace timeout")
		ch <- Timeout
	}
}
