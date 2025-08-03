package k8s

import (
	"context"
	"errors"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func (c *ClusterConnection) CreateSecretObject(
	name string,
	namespace string,
	labels map[string]string,
	annotations map[string]string,
	envVars map[string]string,
) *v1.Secret {
	newMap := map[string][]byte{}
	for k, v := range envVars {
		newMap[k] = []byte(v)
	}
	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Labels:      labels,
			Annotations: annotations,
		},
		Data: newMap,
	}
	return secret
}

func (c *ClusterConnection) ApplySecret(s *v1.Secret) error {
	secret, _ := c.Clientset.CoreV1().Secrets(s.Namespace).Get(context.Background(), s.Name, metav1.GetOptions{})
	if secret.Name != s.Name {
		_, err := c.Clientset.CoreV1().Secrets(s.Namespace).Create(context.Background(), s, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	} else {
		_, err := c.Clientset.CoreV1().Secrets(s.Namespace).Update(context.Background(), s, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *ClusterConnection) DeleteSecret(sName string, sNamespace string) error {
	secret, _ := c.Clientset.CoreV1().Secrets(sNamespace).Get(context.Background(), sName, metav1.GetOptions{})
	if secret.Name != sName {
		return nil
	}
	err := c.Clientset.CoreV1().Secrets(sNamespace).Delete(context.Background(), sName, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (c *ClusterConnection) PatchSecret(sName, sNamespace string, patch []byte) error {
	secret, _ := c.Clientset.CoreV1().Secrets(sNamespace).Get(context.Background(), sName, metav1.GetOptions{})
	if secret.Name != sName {
		return errors.New("secret not found, nothing to patch")
	}
	_, err := c.Clientset.CoreV1().Secrets(sNamespace).Patch(context.Background(), sName, types.StrategicMergePatchType, patch, metav1.PatchOptions{})
	if err != nil {
		return err
	}
	return nil
}
