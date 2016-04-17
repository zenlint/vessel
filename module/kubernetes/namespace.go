package kubernetes

import (
	// "encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/containerops/vessel/models"
	// "k8s.io/kubernetes/pkg/api"
	// "k8s.io/kubernetes/pkg/util/intstr"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/watch"
)

func CreateNamespace(pipelineVersion *models.PipelineVersion) error {
	piplineMetadata := pipelineVersion.MetaData
	// pipelineStageSpecs := pipelineVersion.StageSpecs
	// Going to support create namespace after we have namespace watch lib
	/*	_, err := CLIENT.Namespaces().Get(piplineMetadata.Namespace)
		if err != nil {*/
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
func WatchNamespaceStatus(labelKey string, labelValue string, timeout int, checkOp string) (string, error) {
	if checkOp != string(watch.Deleted) && checkOp != string(watch.Added) {
		fmt.Errorf("Params checkOp err, checkOp: %v", checkOp)
	}

	//opts := api.ListOptions{FieldSelector: fields.Set{"kind": "pod"}.AsSelector()}
	opts := api.ListOptions{LabelSelector: labels.Set{labelKey: labelValue}.AsSelector()}
	w, err := CLIENT.Namespaces().Watch(opts)
	if err != nil {
		fmt.Errorf("Get watch interface err")
		return "", err
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
			// Pod have phase, so we have to wait for the phase change to the right status when added
			if string(event.Type) == checkOp {
				fmt.Println(event.Object.(*api.Namespace).Status.Phase)

				if (checkOp == string(watch.Deleted)) || ((checkOp != string(watch.Deleted)) &&
					(event.Object.(*api.Namespace).Status.Phase == "Active")) {
					return "OK", nil
				}
			}

		case <-t.C:
			return "TIMEOUT", nil
		}
	}
}
