package deployment

import (
	"errors"
	"fmt"
	"strings"

	"github.com/muplat/muplat-backend/pkg/k8s"
	"github.com/muplat/muplat-backend/repositories"
)

func CreatePostgresDeployment(
	deploymentName,
	projectName,
	deploymentType,
	username string,
	pc *PostgresConfig,
	db *repositories.Database,
	cc *k8s.ClusterConnection,
) error {
	if pc == nil {
		return errors.New("postgres config was not provided")
	}

	var database string
	nameSuffix := k8s.GetNameSuffix(fmt.Sprintf("%s%s", deploymentName, projectName))
	resourceName := strings.ToLower(fmt.Sprintf("%s-%s", deploymentName, nameSuffix))
	resourceLabels := map[string]string{
		"name":         deploymentName,
		"created":      username,
		"project-name": projectName,
	}

	if pc.Database == nil {
		database = "app"
	} else {
		database = *pc.Database
	}

	p, err := db.GetProjectByName(projectName)
	if err != nil {
		return err
	}
	projectNamespace := p.Namespace

	u, err := db.GetUserByUsername(username)
	if err != nil {
		return err
	}
	if p.Owner != username && !u.Admin {
		return errors.New("user lacks permissions to create deployment")
	}

	postgresClusterObject := cc.CreatePostgresClusterObject(
		resourceName,
		projectNamespace,
		database,
		resourceLabels,
		map[string]string{},
		pc.DiskSize,
	)
	err = cc.ApplyPostgresCluster(postgresClusterObject)
	if err != nil {
		return err
	}

	err = db.SaveDeployment(deploymentName, projectName, deploymentType)
	if err != nil {
		return err
	}
	err = db.SavePostgresConfig(
		deploymentName,
		projectName,
		fmt.Sprintf("%dGi", pc.DiskSize),
		fmt.Sprintf("%s-rw:5432", resourceName),
		database,
		fmt.Sprintf("%s-app", resourceName),
	)
	if err != nil {
		return err
	}

	return nil
}

func DeletePostgresDeployment(
	deploymentName string,
	projectName string,
	username string,
	db *repositories.Database,
	cc *k8s.ClusterConnection,
) error {
	nameSuffix := k8s.GetNameSuffix(fmt.Sprintf("%s%s", deploymentName, projectName))
	resourceName := strings.ToLower(fmt.Sprintf("%s-%s", deploymentName, nameSuffix))

	p, err := db.GetProjectByName(projectName)
	if err != nil {
		return err
	}
	projectNamespace := p.Namespace

	u, err := db.GetUserByUsername(username)
	if err != nil {
		return err
	}
	if p.Owner != username && !u.Admin {
		return errors.New("user lacks permissions to delete deployment")
	}

	err = cc.DeletePostgresCluster(resourceName, projectNamespace)
	if err != nil {
		return err
	}

	err = db.DeletePostgresConfig(deploymentName, projectName)
	if err != nil {
		return err
	}

	err = db.DeleteDeployment(deploymentName, projectName)
	if err != nil {
		return err
	}

	return nil
}
