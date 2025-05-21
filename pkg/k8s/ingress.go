package k8s

import (

	//v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"context"
	"errors"
	"fmt"

	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func (c *ClusterConnection) CreateIngressObject(
	name string,
	namespace string,
	labels map[string]string,
	annotations map[string]string,
	domainName string,
	serviceName string,
	servicePort uint,
) *v1.Ingress {
	pathType := v1.PathTypePrefix
	ingress := &v1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Labels:      labels,
			Annotations: annotations,
		},
		Spec: v1.IngressSpec{
			IngressClassName: &c.IngressClassName,
			TLS: []v1.IngressTLS{
				{
					Hosts:      []string{domainName},
					SecretName: fmt.Sprintf("crt-%s", name),
				},
			},
			Rules: []v1.IngressRule{
				{
					Host: domainName,
					IngressRuleValue: v1.IngressRuleValue{
						HTTP: &v1.HTTPIngressRuleValue{
							Paths: []v1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &pathType,
									Backend: v1.IngressBackend{
										Service: &v1.IngressServiceBackend{
											Name: serviceName,
											Port: v1.ServiceBackendPort{
												Number: int32(servicePort),
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	return ingress
}

func (c *ClusterConnection) ApplyIngress(i *v1.Ingress) error {
	ingress, _ := c.Clientset.NetworkingV1().Ingresses(i.Namespace).Get(context.Background(), i.Name, metav1.GetOptions{})
	if ingress.Name != i.Name {
		_, err := c.Clientset.NetworkingV1().Ingresses(i.Namespace).Create(context.Background(), i, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	} else {
		_, err := c.Clientset.NetworkingV1().Ingresses(i.Namespace).Update(context.Background(), i, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *ClusterConnection) DeleteIngress(iName string, iNamespace string) error {
	ingress, _ := c.Clientset.NetworkingV1().Ingresses(iNamespace).Get(context.Background(), iName, metav1.GetOptions{})
	if ingress.Name != iName {
		return nil
	}
	err := c.Clientset.NetworkingV1().Ingresses(iNamespace).Delete(context.Background(), iName, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (c *ClusterConnection) PatchIngress(iName, iNamespace string, patch []byte) error {
	ingress, _ := c.Clientset.NetworkingV1().Ingresses(iNamespace).Get(context.Background(), iName, metav1.GetOptions{})
	if ingress.Name != iName {
		return errors.New("ingress not found, nothing to patch")
	}
	_, err := c.Clientset.NetworkingV1().Ingresses(iNamespace).Patch(context.Background(), iName, types.StrategicMergePatchType, patch, metav1.PatchOptions{})
	if err != nil {
		return err
	}
	return nil
}
