package deployment

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/muplat/muplat-backend/repositories"
)

func GetDeployment(deploymentName, projectName, callerUsername string, db *repositories.Database) (*gin.H, error) {
	u, err := db.GetUserByUsername(callerUsername)
	if err != nil {
		return nil, err
	}

	p, err := db.GetProjectByName(projectName)
	if err != nil {
		return nil, err
	}

	if !u.Admin && u.Username != p.Owner {
		return nil, errors.New("you are not allowed to view this deployment")
	}

	d, err := db.GetDeployment(deploymentName, projectName)
	if err != nil {
		return nil, err
	}
	config := map[string]interface{}{}
	switch d.Type {
	case "app":
		ac, err := db.GetAppConfig(deploymentName, projectName)
		if err != nil {
			return nil, err
		}
		config["repository"] = ac.Repository
		config["tag"] = ac.Tag
		config["internalUrl"] = ac.InternalUrl
		config["tier"] = ac.Tier
		if d.AppConfig.ExternalUrl != "" {
			config["externalUrl"] = ac.ExternalUrl
		}
	case "postgres":
		pc, err := db.GetPostgresConfig(deploymentName, projectName)
		if err != nil {
			return nil, err
		}
		config["diskSize"] = pc.DiskSize
		config["internalEndpoint"] = pc.InternalEndpoint
		config["database"] = pc.Database
	}
	response := &gin.H{
		"name":   d.Name,
		"type":   d.Type,
		"config": config,
	}
	return response, nil
}

func GetDeployments(projectName, callerUsername string, db *repositories.Database) (*gin.H, error) {
	u, err := db.GetUserByUsername(callerUsername)
	if err != nil {
		return nil, err
	}

	p, err := db.GetProjectByName(projectName)
	if err != nil {
		return nil, err
	}

	if !u.Admin && u.Username != p.Owner {
		return nil, errors.New("you are not allowed to view deployments of this project")
	}

	dbDeployments, err := db.GetDeploymentsByProject(projectName)
	if err != nil {
		return nil, err
	}

	responseDeployments := []DeploymentResponse{}
	for _, d := range dbDeployments {
		responseDeployments = append(responseDeployments, DeploymentResponse{d.Name, d.Type, d.CreatedAt})
	}

	response := &gin.H{
		"deployments": responseDeployments,
	}

	return response, err
}
