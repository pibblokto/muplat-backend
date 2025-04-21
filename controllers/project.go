package controllers

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	"crypto/sha1"

	"github.com/gin-gonic/gin"
	"github.com/muplat/muplat-backend/models"
	"github.com/muplat/muplat-backend/pkg/jwt"
	"github.com/muplat/muplat-backend/pkg/k8s"
	"github.com/muplat/muplat-backend/pkg/setup"
)

var cfg setup.MuplatCfg = setup.LoadConfig()

type ProjectInput struct {
	Name string `json:"name" binding:"required"`
}

func CreateProject(c *gin.Context) {
	var input ProjectInput
	var ingressNginxNamespace string = cfg.IngressNginxNamespace
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input.Name = strings.TrimSpace(input.Name)

	username, err := jwt.ExtractTokenUsername(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	hasher := sha1.New()
	hasher.Write(fmt.Appendf(nil, "%s%s", input.Name, username))
	nameSuffix := hex.EncodeToString(hasher.Sum(nil))
	nameSuffix = nameSuffix[:7]

	namespaceName := strings.ToLower(fmt.Sprintf("%s-%s", input.Name, nameSuffix))
	networkPolicyName := strings.ToLower(fmt.Sprintf("%s-%s", input.Name, nameSuffix))

	namespaceLabels := map[string]string{
		"name":         namespaceName,
		"owner":        username,
		"project-name": input.Name,
	}
	networkPolicyLabels := map[string]string{
		"name":         networkPolicyName,
		"owner":        username,
		"project-name": input.Name,
	}

	clientset, err := k8s.ConnectCluster()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	namespaceObject := k8s.CreateNamespaceObject(namespaceName, namespaceLabels, map[string]string{})
	err = k8s.CreateNamespace(clientset, namespaceObject)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	namespacePolicyObject := k8s.CreateNetworkPolicyObject(
		networkPolicyName,
		namespaceName,
		networkPolicyLabels,
		map[string]string{},
		ingressNginxNamespace,
	)
	err = k8s.CreateNetworkPolicy(clientset, namespacePolicyObject)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	owner, err := jwt.ExtractTokenUsername(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newProject := models.Project{
		Name:          input.Name,
		Owner:         owner,
		Namespace:     namespaceName,
		NetworkPolicy: networkPolicyName,
		ProjectHash:   nameSuffix,
	}
	_, err = newProject.SaveProject()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "project " + input.Name + " was created"})
}

func DeleteProject(c *gin.Context) {
	var input ProjectInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	p := models.Project{}
	err := p.GetPorjectByName(input.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username, err := jwt.ExtractTokenUsername(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	u := models.User{}
	err = u.GetUserByUsername(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if p.Owner != username && !u.Admin {
		c.JSON(http.StatusForbidden, gin.H{"error": "only admin or owner of the project can delete it"})
		return
	}

	clientset, err := k8s.ConnectCluster()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = k8s.DeleteNetworkPolicy(clientset, p.NetworkPolicy, p.Namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = k8s.DeleteNaspace(clientset, p.Namespace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = p.DeleteProject()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "project " + input.Name + " was deleted"})
}
