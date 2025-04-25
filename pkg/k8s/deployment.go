package k8s

import (
	"context"
	"fmt"

	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type appTier struct {
	Replicas int32
}

var (
	appTiers = map[string]appTier{
		"dev": {
			Replicas: 1,
		},
		"mid": {
			Replicas: 2,
		},
		"pro": {
			Replicas: 5,
		},
	}
)

func CreateDeploymentObject(
	name string,
	namespace string,
	labels map[string]string,
	annotations map[string]string,
	tier string,
	repository string,
	tag string,
	port uint,
	envSecretName string,
) *v1.Deployment {
	replicas := appTiers[tier].Replicas
	deployment := &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Labels:      labels,
			Annotations: annotations,
		},
		Spec: v1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"name": labels["name"],
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"name": labels["name"],
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  name,
							Image: fmt.Sprintf("%s:%s", repository, tag),
							Ports: []corev1.ContainerPort{
								{
									Name:          GetPortName(name),
									ContainerPort: int32(port),
								},
							},
						},
					},
				},
			},
		},
	}

	if envSecretName != "" {
		deployment.Spec.Template.Spec.Containers[0].EnvFrom = []corev1.EnvFromSource{
			{
				SecretRef: &corev1.SecretEnvSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: envSecretName,
					},
				},
			},
		}
	}
	return deployment
}

func ApplyDeployment(clientset *kubernetes.Clientset, d *v1.Deployment) error {
	deployment, _ := clientset.AppsV1().Deployments(d.Namespace).Get(context.Background(), d.Name, metav1.GetOptions{})
	if deployment.Name != d.Name {
		_, err := clientset.AppsV1().Deployments(d.Namespace).Create(context.Background(), d, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	} else {
		_, err := clientset.AppsV1().Deployments(d.Namespace).Update(context.Background(), d, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

func DeleteDeployment(clientset *kubernetes.Clientset, dName string, dNamespace string) error {
	deployment, _ := clientset.AppsV1().Deployments(dNamespace).Get(context.Background(), dName, metav1.GetOptions{})
	if deployment.Name != dName {
		return nil
	}
	err := clientset.AppsV1().Deployments(dNamespace).Delete(context.Background(), dName, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}
