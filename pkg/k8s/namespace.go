package k8s

import (

	//v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func CreateNamespaceObject(
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

func CreateNamespace(clientset *kubernetes.Clientset, ns *v1.Namespace) error {
	namespace, _ := clientset.CoreV1().Namespaces().Get(context.Background(), ns.Name, metav1.GetOptions{})
	if namespace.Name != ns.Name {
		_, err := clientset.CoreV1().Namespaces().Create(context.Background(), ns, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}
