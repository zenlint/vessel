package kubernetes

import (
	// "encoding/json"
	"fmt"

	"github.com/containerops/vessel/models"
	// "k8s.io/kubernetes/pkg/api"
	// "k8s.io/kubernetes/pkg/util/intstr"
)

func StartPipeline(pipelineVersion *models.PipelineVersion) error {
	piplineMetadata := pipelineVersion.MetaData
	if _, err := CLIENT.Namespaces().Get(piplineMetadata.Namespace); err != nil {
		fmt.Println("111111111111111")
		if err := CreateNamespace(pipelineVersion); err != nil {
			return err
		}
	}
	fmt.Println("222222222222222222222222")

	/*if status, err := WatchNamespaceStatus("app", piplineMetadata.Name, 30, Added); err != nil || status != "OK" {
		// if status != "OK" {
		return err
		// }
	}
	*/
	if err := CreateRC(pipelineVersion); err != nil {
		return err
	}

	if err := CreateService(pipelineVersion); err != nil {
		return err
	}
	//createrc && createservice
	return nil
}
