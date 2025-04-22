package k8s

import (
	"context"
	b64 "encoding/base64"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func CreateSecretObject(
	name string,
	namespace string,
	labels map[string]string,
	annotations map[string]string,
	envVars map[string]string,
) *v1.Secret {
	encodedMap := map[string][]byte{}
	for k, v := range envVars {
		encodedMap[k] = []byte(b64.StdEncoding.EncodeToString([]byte(v)))
	}
	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Labels:      labels,
			Annotations: annotations,
		},
		Data: encodedMap,
	}
	return secret
}

func ApplySecret(clientset *kubernetes.Clientset, s *v1.Secret) error {
	secret, _ := clientset.CoreV1().Secrets(s.Namespace).Get(context.Background(), s.Name, metav1.GetOptions{})
	if secret.Name != s.Name {
		_, err := clientset.CoreV1().Secrets(s.Namespace).Create(context.Background(), s, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	} else {
		_, err := clientset.CoreV1().Secrets(s.Namespace).Update(context.Background(), s, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

func DeleteSecret(clientset *kubernetes.Clientset, sName string, sNamespace string) error {
	secret, _ := clientset.NetworkingV1().Ingresses(sNamespace).Get(context.Background(), sName, metav1.GetOptions{})
	if secret.Name != sName {
		return nil
	}
	err := clientset.CoreV1().Secrets(sNamespace).Delete(context.Background(), sName, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}
