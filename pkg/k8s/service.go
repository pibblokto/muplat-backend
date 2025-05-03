package k8s

import (
	"context"
	"errors"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func (c *ClusterConnection) CreateServiceObject(
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

func (c *ClusterConnection) ApplyService(s *v1.Service) error {
	service, _ := c.Clientset.CoreV1().Services(s.Namespace).Get(context.Background(), s.Name, metav1.GetOptions{})
	if service.Name != s.Name {
		_, err := c.Clientset.CoreV1().Services(s.Namespace).Create(context.Background(), s, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	} else {
		_, err := c.Clientset.CoreV1().Services(s.Namespace).Update(context.Background(), s, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *ClusterConnection) DeleteService(sName string, sNamespace string) error {
	service, _ := c.Clientset.CoreV1().Services(sNamespace).Get(context.Background(), sName, metav1.GetOptions{})
	if service.Name != sName {
		return nil
	}
	err := c.Clientset.CoreV1().Services(sNamespace).Delete(context.Background(), sName, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (c *ClusterConnection) GetNginxControllerIp() (string, error) {
	serviceList, err := c.Clientset.CoreV1().Services(c.IngressNamespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return "", nil
	}
	for _, s := range serviceList.Items {
		if s.Spec.Type == v1.ServiceTypeLoadBalancer {
			return s.Spec.LoadBalancerIP, nil
		}
	}
	return "", errors.New("nginx controller has no LoadBalancer services in its namespace")
}
