package k8s

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

func CreateServiceObject(
	name string,
	namespace string,
	labels map[string]string,
	annotations map[string]string,
	deploymentNameLabel string,
	port uint,
) *v1.Service {
	return &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Labels:      labels,
			Annotations: annotations,
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Port: int32(port),
					TargetPort: intstr.IntOrString{
						Type:   0,
						IntVal: int32(port),
					},
				},
			},
			Selector: map[string]string{
				"name": deploymentNameLabel,
			},
			Type: v1.ServiceTypeClusterIP,
		},
	}
}

func ApplyService(clientset *kubernetes.Clientset, s *v1.Service) error {
	service, _ := clientset.CoreV1().Services(s.Namespace).Get(context.Background(), s.Name, metav1.GetOptions{})
	if service.Name != s.Name {
		_, err := clientset.CoreV1().Services(s.Namespace).Create(context.Background(), s, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	} else {
		_, err := clientset.CoreV1().Services(s.Namespace).Update(context.Background(), s, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}
