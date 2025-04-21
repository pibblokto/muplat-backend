package controllers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/muplat/muplat-backend/models"
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

	nameSuffix := k8s.GetNameSuffix(fmt.Sprintf("%s%s%s%s", input.Name, input.ProjectName, username))
	deploymentName := strings.ToLower(fmt.Sprintf("%s-%s", input.Name, nameSuffix))
	deploymentNamespace := p.Namespace
	secretName := strings.ToLower(fmt.Sprintf("%s-%s", input.Name, nameSuffix))

	deploymentLabels := map[string]string{
		"name":         deploymentName,
		"owner":        username,
		"project-name": p.Name,
	}
	secretLabels := map[string]string{
		"name":         secretName,
		"owner":        username,
		"project-name": p.Name,
	}

	clientset, err := k8s.ConnectCluster()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	switch input.Type {
	case TypeApp:
		fmt.Println("to be added")
	}
}
