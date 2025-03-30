package k8s

import (

	//v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/muplat/muplat-backend/pkg/bootstrap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var appConf bootstrap.AppCfg = bootstrap.LoadConfig()

func ConnectCluster() (*kubernetes.Clientset, error) {

	config, err := clientcmd.BuildConfigFromFlags("", appConf.KubeconfigPath)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func GenerateNamespaceObject(name string, labels map[string]string, annotations map[string]string) *corev1.Namespace {
	return &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Labels:      labels,
			Annotations: annotations,
		},
	}
}
