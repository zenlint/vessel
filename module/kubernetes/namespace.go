package kubernetes

import (
	"github.com/containerops/vessel/models"
	"k8s.io/kubernetes/pkg/api"
)

func createNamespace(stage *models.Stage) error {
	if err := getClient(); err != nil {
		return err
	}
	k8sNamespace := k8sClient.Namespaces()
	namespaceLock.RLock()
	if _, err := k8sNamespace.Get(stage.Namespace); err != nil {
		namespaceLock.RUnlock()
		namespaceLock.Lock()
		if _, err := k8sNamespace.Get(stage.Namespace); err != nil {
			namespaceObj := &api.Namespace{
				ObjectMeta: api.ObjectMeta{
					Name:   stage.Namespace,
					Labels: map[string]string{},
				},
			}
			namespaceObj.SetLabels(map[string]string{models.LabelKey: stage.PipelineName})

			if _, err := k8sNamespace.Create(namespaceObj); err != nil {
				namespaceLock.Unlock()
				return err
			}
		}
		namespaceLock.Unlock()
	} else {
		namespaceLock.RUnlock()
	}
	return nil
}

func deleteNamespace(stage *models.Stage) error {
	if err := getClient(); err != nil {
		return err
	}
	k8sNamespace := k8sClient.Namespaces()
	namespaceLock.RLock()
	if _, err := k8sNamespace.Get(stage.Namespace); err == nil {
		namespaceLock.RUnlock()
		namespaceLock.Lock()
		if _, err := k8sNamespace.Get(stage.Namespace); err == nil {
			if err := k8sNamespace.Delete(stage.Namespace); err != nil {
				namespaceLock.Unlock()
				return err
			}
		}
		namespaceLock.Unlock()
	} else {
		namespaceLock.RUnlock()
	}
	return nil
}
