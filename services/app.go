package services

import (
	"fmt"
	"strings"

	"github.com/muplat/muplat-backend/pkg/k8s"
	"github.com/muplat/muplat-backend/repositories"
	"k8s.io/client-go/kubernetes"
)

type AppService struct {
	Db  *repositories.DatabaseConfig
	K8s *k8s.ClusterConfig
}

func CreateAppDeployment(
	clientset *kubernetes.Clientset,
	deploymentName,
	projectName,
	projectNamespace,
	owner string,
	ac AppConfig,
) (string, error) {
	nameSuffix := k8s.GetNameSuffix(fmt.Sprintf("%s%s", deploymentName, projectName))
	resourceName := strings.ToLower(fmt.Sprintf("%s-%s", deploymentName, nameSuffix))
	resourceLabels := map[string]string{
		"name":         deploymentName,
		"created":      owner,
		"project-name": projectName,
	}

	var secretName string
	if ac.EnvVars != nil {
		secretObject := k8s.CreateSecretObject(
			resourceName,
			projectNamespace,
			resourceLabels,
			map[string]string{},
			ac.EnvVars,
		)
		err := k8s.ApplySecret(clientset, secretObject)
		if err != nil {
			return "", err
		}
		secretName = resourceName
	} else {
		secretName = ""
	}

	// Deployment
	deploymentObject := k8s.CreateDeploymentObject(
		resourceName,
		projectNamespace,
		resourceLabels,
		map[string]string{},
		string(ac.Tier),
		ac.Repository,
		ac.Tag,
		ac.Port,
		secretName,
	)
	err := k8s.ApplyDeployment(clientset, deploymentObject)
	if err != nil {
		deleteError := DeleteAppDeployment(
			clientset,
			deploymentName,
			projectName,
			projectNamespace,
		)
		if deleteError != nil {
			return "", fmt.Errorf("failed to create app %s and delete orphan resources", deploymentName)
		}
		return "", err
	}

	// Service
	serviceObject := k8s.CreateServiceObject(
		resourceName,
		projectNamespace,
		resourceLabels,
		map[string]string{},
		resourceLabels["name"],
		ac.Port,
	)
	err = k8s.ApplyService(clientset, serviceObject)

	if err != nil {
		deleteError := DeleteAppDeployment(
			clientset,
			deploymentName,
			projectName,
			projectNamespace,
		)
		if deleteError != nil {
			return "", fmt.Errorf("failed to create app %s and delete orphan resources", deploymentName)
		}
		return "", err
	}

	// Ingress
	if *ac.External {
		ingressObject := k8s.CreateIngressObject(
			resourceName,
			projectNamespace,
			resourceLabels,
			map[string]string{},
			cfg.IngressClassName,
			ac.DomainName,
			resourceName,
			ac.Port,
		)
		err = k8s.ApplyIngress(clientset, ingressObject)
		if err != nil {
			deleteError := DeleteAppDeployment(
				clientset,
				deploymentName,
				projectName,
				projectNamespace,
			)
			if deleteError != nil {
				return "", fmt.Errorf("failed to create app %s and delete orphan resources", deploymentName)
			}
			return "", err
		}
	}
	return resourceName, nil
}

func DeleteAppDeployment(
	clientset *kubernetes.Clientset,
	deploymentName string,
	projectName string,
	projectNamespace string,
) error {
	nameSuffix := k8s.GetNameSuffix(fmt.Sprintf("%s%s", deploymentName, projectName))
	resourceName := strings.ToLower(fmt.Sprintf("%s-%s", deploymentName, nameSuffix))

	err := k8s.DeleteIngress(clientset, resourceName, projectNamespace)
	if err != nil {
		return err
	}

	err = k8s.DeleteService(clientset, resourceName, projectNamespace)
	if err != nil {
		return err
	}

	err = k8s.DeleteDeployment(clientset, resourceName, projectNamespace)
	if err != nil {
		return err
	}

	err = k8s.DeleteSecret(clientset, resourceName, projectNamespace)
	if err != nil {
		return err
	}

	return nil
}
