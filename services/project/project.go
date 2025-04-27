package project

import (
	"errors"
	"fmt"
	"strings"

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
	p, err := db.GetPorjectByName(projectName)
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
