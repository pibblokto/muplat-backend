package controllers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/muplat/muplat-backend/pkg/jwt"
	"github.com/muplat/muplat-backend/pkg/k8s"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type RequestNS struct {
	Name        string            `json:"name" binding:"required"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
}

func CreateNamespace(c *gin.Context) {

	var ns RequestNS
	if err := c.ShouldBindJSON(&ns); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clientset, err := k8s.ConnectCluster()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	namespace, _ := clientset.CoreV1().Namespaces().Get(context.TODO(), ns.Name, metav1.GetOptions{})
	if namespace.Name != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Namespace " + namespace.Name + " already exists"})
		return
	}

	namespaceObject := k8s.CreateNamespaceObject(ns.Name, ns.Labels, ns.Annotations)

	_, err = clientset.CoreV1().Namespaces().Create(context.TODO(), namespaceObject, metav1.CreateOptions{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "namespace " + ns.Name + " was created"})
}

type ProjectInput struct {
	Name string `json:"username" binding:"required"`
}

func CreateProject(c *gin.Context) {
	var project ProjectInput
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	username, err := jwt.ExtractTokenUsername(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	namespaceName := fmt.Sprintf("%s-%s", username, project.Name)
	networkPolicyName := fmt.Sprintf("%s-%s", username, project.Name)

	namespaceLabels := map[string]string{
		"name":  namespaceName,
		"owner": username,
	}
	networkPolicyLabels := map[string]string{
		"name":  networkPolicyName,
		"owner": username,
	}
}
