package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/muplat/muplat-backend/services/deployment"
)

func (h *HttpHandler) CreateDeployment(c *gin.Context) {
	var input CreateDeploymentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username, err := h.Jwt.ExtractTokenUsername(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var ipForRecord string
	switch input.Type {
	case deployment.TypeApp:
		ipForRecord, err = deployment.CreateAppDeployment(
			input.Name,
			input.ProjectName,
			string(input.Type),
			username,
			input.AppConfing,
			h.Db,
			h.ClusterConn,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	case deployment.TypePostgres:
		err := deployment.CreatePostgresDeployment(
			input.Name,
			input.ProjectName,
			string(input.Type),
			username,
			input.PostgresConfig,
			h.Db,
			h.ClusterConn,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "deployment of this type doesn't exist"})
		return
	}

	response := gin.H{"message": "deployment " + input.Name + " was created"}
	if ipForRecord != "" {
		response["ipForRecord"] = ipForRecord
	}
	c.JSON(http.StatusOK, response)
}

func (h *HttpHandler) DeleteDeployment(c *gin.Context) {
	var input DeleteDeploymentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username, err := h.Jwt.ExtractTokenUsername(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	d, err := h.Db.GetDeployment(input.Name, input.ProjectName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	deploymentType := deployment.DeploymentType(d.Type)

	switch deploymentType {
	case deployment.TypeApp:
		err := deployment.DeleteAppDeployment(
			input.Name,
			input.ProjectName,
			username,
			h.Db,
			h.ClusterConn,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	case deployment.TypePostgres:
		err := deployment.DeletePostgresDeployment(
			input.Name,
			input.ProjectName,
			username,
			h.Db,
			h.ClusterConn,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "deployment " + input.Name + " was deleted"})
}

func (h *HttpHandler) GetDeployment(c *gin.Context) {
	projectName := c.Param("project")
	deploymentName := c.Param("deployment")
	username, err := h.Jwt.ExtractTokenUsername(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response, err := deployment.GetDeployment(deploymentName, projectName, username, h.Db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *HttpHandler) GetDeployments(c *gin.Context) {
	projectName := c.Param("project")
	username, err := h.Jwt.ExtractTokenUsername(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response, err := deployment.GetDeployments(projectName, username, h.Db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *HttpHandler) ReissueCertificate(c *gin.Context) {
	var input ReissueCertificateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username, err := h.Jwt.ExtractTokenUsername(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = deployment.ReissueCertificate(input.Name, input.ProjectName, username, h.Db, h.ClusterConn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "reissuing started"})
}
