package kubernetes

/*
// Delete takes the name of the pod, and returns an error if one occurs
func (c *pods) Delete(name string, options *api.DeleteOptions) error {
	return c.r.Delete().Namespace(c.ns).Resource("pods").Name(name).Body(options).Do().Error()
}

*/

/*func DeleteK8SResource(pipeline PiplelineInterface) error {
	var pvm piplineMetadata
	var pvs piplineSpec
	err := split(pipeline, &pvm, &pvs)
	if err != nil {
		return err
	}

	namespace, name := getNamespaceAndName(pvm, pvs)
	if err = CLIENT.ReplicationControllers(namespace).Delete(name); err != nil {
		return err
	}
	if err = CLIENT.Pods(namespace).Delete(name, nil); err != nil {
		return err
	}
	if err = CLIENT.Services(namespace).Delete(name); err != nil {
		return err
	}

	return nil
}

func getNamespaceAndName(plm piplineMetadata, pls piplineSpec) (string, string) {
	return plm.namespace, pls.name
}
*/
