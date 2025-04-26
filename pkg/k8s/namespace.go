package k8s

import (

	//v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *ClusterConfig) CreateNamespaceObject(
	name string,
	labels map[string]string,
	annotations map[string]string,
) *v1.Namespace {
	return &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Labels:      labels,
			Annotations: annotations,
		},
	}
}

func (c *ClusterConfig) ApplyNamespace(ns *v1.Namespace) error {
	namespace, _ := c.Clientset.CoreV1().Namespaces().Get(context.Background(), ns.Name, metav1.GetOptions{})
	if namespace.Name != ns.Name {
		_, err := c.Clientset.CoreV1().Namespaces().Create(context.Background(), ns, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	} else {
		_, err := c.Clientset.CoreV1().Namespaces().Update(context.Background(), ns, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *ClusterConfig) DeleteNamespace(ns string) error {
	namespace, _ := c.Clientset.CoreV1().Namespaces().Get(context.Background(), ns, metav1.GetOptions{})
	if namespace.Name != ns {
		return nil
	}
	err := c.Clientset.CoreV1().Namespaces().Delete(context.Background(), ns, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}
