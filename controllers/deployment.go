package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/muplat/muplat-backend/models"
	"github.com/muplat/muplat-backend/pkg/deployments"
	"github.com/muplat/muplat-backend/pkg/jwt"
	"github.com/muplat/muplat-backend/pkg/k8s"
)

func CreateDeployment(c *gin.Context) {
	var input DeploymentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username, err := jwt.ExtractTokenUsername(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	p := models.Project{}
	err = p.GetPorjectByName(input.ProjectName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := models.User{}
	err = u.GetUserByUsername(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if p.Owner != u.Username && !u.Admin {
		c.JSON(http.StatusForbidden, gin.H{"error": "you have to be an admin or an owner of the project"})
		return
	}

	clientset, err := k8s.ConnectCluster()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	switch input.Type {
	case deployments.TypeApp:
		if input.AppConfing == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "app config was not provided"})
			return
		}
		resourceName, err := deployments.CreateAppDeployment(
			clientset,
			input.Name,
			input.ProjectName,
			p.Namespace,
			username,
			*input.AppConfing,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		newDeployment := models.Deployment{
			Name:        input.Name,
			ProjectName: input.ProjectName,
			Type:        string(input.Type),
		}
		_, err = newDeployment.SaveDeployment()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		newAppConfig := models.AppConfig{
			DeploymentName: input.Name,
			ProjectName:    input.ProjectName,
			Repository:     input.AppConfing.Repository,
			Tag:            input.AppConfing.Tag,
			InternalUrl:    fmt.Sprintf("http://%s", resourceName),
			Tier:           string(input.AppConfing.Tier),
			Port:           input.AppConfing.Port,
			EnvVarsSecret:  resourceName,
		}
		if *input.AppConfing.External {
			newAppConfig.ExternalUrl = fmt.Sprintf("http://%s", input.AppConfing.DomainName)
		}
		_, err = newAppConfig.SaveAppConfig()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "deployment " + input.Name + " was created"})
}

func DeleteDeployment(c *gin.Context) {
	var input DeploymentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username, err := jwt.ExtractTokenUsername(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	p := models.Project{}
	err = p.GetPorjectByName(input.ProjectName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u := models.User{}
	err = u.GetUserByUsername(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if p.Owner != u.Username && !u.Admin {
		c.JSON(http.StatusForbidden, gin.H{"error": "you have to be an admin or an owner of the project"})
		return
	}

	clientset, err := k8s.ConnectCluster()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	switch input.Type {
	case deployments.TypeApp:
		err := deployments.DeleteAppDeployment(
			clientset,
			input.Name,
			input.ProjectName,
			p.Namespace,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ac := models.AppConfig{
			DeploymentName: input.Name,
			ProjectName:    input.ProjectName,
		}
		err = ac.DeleteAppConfig()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		d := models.Deployment{
			Name:        input.Name,
			ProjectName: input.ProjectName,
		}
		err = d.DeleteDeployment()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "deployment " + input.Name + " was deleted"})
}
