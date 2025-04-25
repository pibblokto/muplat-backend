package deployments

import (
	"fmt"
	"strings"

	"github.com/muplat/muplat-backend/pkg/k8s"
	"k8s.io/client-go/dynamic"
)

func CreatePostgresDeployment(
	client *dynamic.DynamicClient,
	deploymentName,
	projectName,
	projectNamespace,
	owner string,
	pc PostgresConfig,
) (string, error) {
	var database string
	nameSuffix := k8s.GetNameSuffix(fmt.Sprintf("%s%s", deploymentName, projectName))
	resourceName := strings.ToLower(fmt.Sprintf("%s-%s", deploymentName, nameSuffix))
	resourceLabels := map[string]string{
		"name":         deploymentName,
		"created":      owner,
		"project-name": projectName,
	}
	if pc.Database == nil {
		database = "app"
	} else {
		database = *pc.Database
	}

	postgresClusterObject := k8s.CreatePostgresClusterObject(
		resourceName,
		projectNamespace,
		database,
		resourceLabels,
		map[string]string{},
		pc.DiskSize,
	)
	err := k8s.ApplyPostgresCluster(client, postgresClusterObject)
	if err != nil {
		return "", err
	}

	return resourceName, nil
}

func DeletePostgresDeployment(
	client *dynamic.DynamicClient,
	deploymentName string,
	projectName string,
	projectNamespace string,
) error {
	nameSuffix := k8s.GetNameSuffix(fmt.Sprintf("%s%s", deploymentName, projectName))
	resourceName := strings.ToLower(fmt.Sprintf("%s-%s", deploymentName, nameSuffix))

	err := k8s.DeletePostgresCluster(client, resourceName, projectNamespace)
	if err != nil {
		return err
	}

	return nil
}
