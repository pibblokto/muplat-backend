package k8s

import (
	"context"
	"errors"
	"fmt"

	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func (c *ClusterConnection) CreateDeploymentObject(
	name string,
	namespace string,
	labels map[string]string,
	annotations map[string]string,
	repository string,
	tag string,
	port uint,
	envSecretName string,
	replicas int32,
) *v1.Deployment {
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

func (c *ClusterConnection) ApplyDeployment(d *v1.Deployment) error {
	deployment, _ := c.Clientset.AppsV1().Deployments(d.Namespace).Get(context.Background(), d.Name, metav1.GetOptions{})
	if deployment.Name != d.Name {
		_, err := c.Clientset.AppsV1().Deployments(d.Namespace).Create(context.Background(), d, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	} else {
		_, err := c.Clientset.AppsV1().Deployments(d.Namespace).Update(context.Background(), d, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *ClusterConnection) DeleteDeployment(dName, dNamespace string) error {
	deployment, _ := c.Clientset.AppsV1().Deployments(dNamespace).Get(context.Background(), dName, metav1.GetOptions{})
	if deployment.Name != dName {
		return nil
	}
	err := c.Clientset.AppsV1().Deployments(dNamespace).Delete(context.Background(), dName, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (c *ClusterConnection) PatchDeployment(dName, dNamespace string, patch []byte) error {
	deployment, _ := c.Clientset.AppsV1().Deployments(dNamespace).Get(context.Background(), dName, metav1.GetOptions{})
	if deployment.Name != dName {
		return errors.New("deployment not found, nothing to patch")
	}
	_, err := c.Clientset.AppsV1().Deployments(dNamespace).Patch(context.Background(), dName, types.StrategicMergePatchType, patch, metav1.PatchOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (c *ClusterConnection) GetDeploymentImage(dName, dNamespace string) (string, error) {
	deployment, _ := c.Clientset.AppsV1().Deployments(dNamespace).Get(context.Background(), dName, metav1.GetOptions{})
	if deployment.Name != dName {
		return "", errors.New("deployment not found")
	}
	return deployment.Spec.Template.Spec.Containers[0].Image, nil
}
