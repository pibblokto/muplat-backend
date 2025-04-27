package k8s

import (

	//v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"

	"github.com/caarlos0/env"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type ConnectionMode string

const (
	ExternalConnection ConnectionMode = "external"
	InternalConnection ConnectionMode = "internal"
)

type ClusterConnection struct {
	Clientset        *kubernetes.Clientset
	Client           *dynamic.DynamicClient
	kubeconfigPath   string         `env:"KUBECONFIG"`
	ingressNamespace string         `env:"INGRESS_NGINX_NAMESPACE" envDefault:"ingress-nginx"`
	ingressClassName string         `env:"INGRESS_CLASS_NAME" envDefault:"nginx"`
	connectionMode   ConnectionMode `env:"CONNECTION_MODE" envDefault:"internal"`
}

func NewClusterConnection() (c *ClusterConnection) {
	var config *rest.Config
	var err error

	err = env.Parse(c)
	if err != nil {
		log.Fatalf("Cluster connection config initialization error: %v", err)
	}

	if c.connectionMode != InternalConnection && c.connectionMode != ExternalConnection {
		log.Fatal("CONNECTION_MODE should be either internal or external")
	}

	if c.connectionMode == "external" && c.kubeconfigPath == "" {
		log.Fatalf("KUBECONFIG is required if CONNECTION_MODE is \"%s\"", c.connectionMode)
	}

	if c.connectionMode == ExternalConnection {
		config, err = clientcmd.BuildConfigFromFlags("", c.kubeconfigPath)
		if err != nil {
			log.Fatalf("Failed to get configuration from kubeconfig: %v", err)
		}

	} else if c.connectionMode == InternalConnection {
		config, err = rest.InClusterConfig()
		if err != nil {
			log.Fatalf("Failed to get configuration for attached ServiceAccount: %v", err)
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to initialize clienset: %v", err)
	}

	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to initialize dynamic client: %v", err)
	}

	c.Clientset = clientset
	c.Client = dynamicClient
	return c
}
