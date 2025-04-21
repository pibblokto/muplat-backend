package k8s

import (

	//v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"context"

	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func CreateIngressObject(
	name string,
	namespace string,
	labels map[string]string,
	annotations map[string]string,
	ingressClassName string,
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
			IngressClassName: &ingressClassName,
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

func ApplyIngress(clientset *kubernetes.Clientset, i *v1.Ingress) error {
	ingress, _ := clientset.NetworkingV1().Ingresses(i.Namespace).Get(context.Background(), i.Name, metav1.GetOptions{})
	if ingress.Name != i.Name {
		_, err := clientset.NetworkingV1().Ingresses(i.Namespace).Create(context.Background(), ingress, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	} else {
		_, err := clientset.NetworkingV1().Ingresses(i.Namespace).Update(context.Background(), ingress, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

// To be changed
func DeleteIngress(clientset *kubernetes.Clientset, npName string, npNamespace string) error {
	networkPolicy, _ := clientset.NetworkingV1().NetworkPolicies(npNamespace).Get(context.Background(), npName, metav1.GetOptions{})
	if networkPolicy.Name != npName {
		return nil
	}
	err := clientset.NetworkingV1().NetworkPolicies(npNamespace).Delete(context.Background(), npName, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}
