package project

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/muplat/muplat-backend/models"
	"github.com/muplat/muplat-backend/pkg/k8s"
	"github.com/muplat/muplat-backend/repositories"
)

func CreateProject(
	projectName string,
	username string,
	db *repositories.Database,
	cc *k8s.ClusterConnection,
) error {

	nameSuffix := k8s.GetNameSuffix(projectName)
	namespaceName := strings.ToLower(fmt.Sprintf("%s-%s", projectName, nameSuffix))
	networkPolicyName := strings.ToLower(fmt.Sprintf("%s-%s", projectName, nameSuffix))

	namespaceLabels := map[string]string{
		"name":         namespaceName,
		"created":      username,
		"project-name": projectName,
	}
	networkPolicyLabels := map[string]string{
		"name":         networkPolicyName,
		"created":      username,
		"project-name": projectName,
	}

	namespaceObject := cc.CreateNamespaceObject(namespaceName, namespaceLabels, map[string]string{})
	err := cc.ApplyNamespace(namespaceObject)
	if err != nil {
		return err
	}

	namespacePolicyObject := cc.CreateNetworkPolicyObject(
		networkPolicyName,
		namespaceName,
		networkPolicyLabels,
		map[string]string{},
	)
	err = cc.ApplyNetworkPolicy(namespacePolicyObject)
	if err != nil {
		return err
	}

	err = db.SaveProject(projectName, username, namespaceName, networkPolicyName)
	if err != nil {
		return err
	}
	return nil
}

func DeleteProject(
	projectName string,
	username string,
	db *repositories.Database,
	cc *k8s.ClusterConnection,
) error {
	p, err := db.GetProjectByName(projectName)
	if err != nil {
		return err
	}

	u, err := db.GetUserByUsername(username)
	if err != nil {
		return err
	}

	if p.Owner != username && !u.Admin {
		return errors.New("user lacks permissions to delete project")
	}

	err = cc.DeleteNetworkPolicy(p.NetworkPolicy, p.Namespace)
	if err != nil {
		return err
	}
	err = cc.DeleteNamespace(p.Namespace)
	if err != nil {
		return err
	}

	err = db.DeleteProject(projectName)
	if err != nil {
		return err
	}
	return nil
}

func GetProject(projectName, callerUsername string, db *repositories.Database) (*gin.H, error) {
	p, err := db.GetProjectByName(projectName)
	if err != nil {
		return nil, err
	}

	u, err := db.GetUserByUsername(callerUsername)
	if err != nil {
		return nil, err
	}

	if !u.Admin && p.Owner != callerUsername {
		return nil, errors.New("you lack permissions to view this project")
	}

	dbDeployments, err := db.GetDeploymentsByProject(projectName)
	if err != nil {
		return nil, err
	}

	responseDeployments := []DeploymentResponse{}
	for _, d := range dbDeployments {
		responseDeployments = append(responseDeployments, DeploymentResponse{d.Name, d.Type, d.CreatedAt})
	}

	projects := &gin.H{
		"name":        p.Name,
		"owner":       p.Owner,
		"createdAt":   p.CreatedAt,
		"deployments": responseDeployments,
	}
	return projects, nil
}

func GetProjects(callerUsername string, db *repositories.Database) (*gin.H, error) {
	u, err := db.GetUserByUsername(callerUsername)
	if err != nil {
		return nil, err
	}

	var dbProjects []*models.Project
	if u.Admin {
		dbProjects, err = db.GetProjects()
		if err != nil {
			return nil, err
		}
	} else {
		dbProjects, err = db.GetProjectsByOwner(u.Username)
		if err != nil {
			return nil, err
		}
	}

	responseProjects := []ProjectResponse{}
	for _, p := range dbProjects {
		responseProjects = append(responseProjects, ProjectResponse{p.Name, p.Owner, p.CreatedAt})
	}
	projects := &gin.H{
		"projects": responseProjects,
	}
	return projects, nil
}
