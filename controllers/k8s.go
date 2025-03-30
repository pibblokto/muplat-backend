package controllers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
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

	namespaceObject := k8s.GenerateNamespaceObject(ns.Name, ns.Labels, ns.Annotations)

	_, err = clientset.CoreV1().Namespaces().Create(context.TODO(), namespaceObject, metav1.CreateOptions{})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "namespace " + ns.Name + " was created"})
}
