package k8s

import (

	//v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/muplat/muplat-backend/pkg/setup"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var cfg setup.MuplatCfg = setup.LoadConfig()

func ConnectCluster() (*kubernetes.Clientset, error) {
	var clientset *kubernetes.Clientset
	var config *rest.Config
	var err error

	if cfg.ConnectionMode == "external" {
		config, err = clientcmd.BuildConfigFromFlags("", cfg.KubeconfigPath)
		if err != nil {
			return nil, err
		}

	} else if cfg.ConnectionMode == "internal" {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	}
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}
