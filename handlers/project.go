package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/muplat/muplat-backend/services/project"
)

type ProjectInput struct {
	Name string `json:"name" binding:"required"`
}

func (h *HttpHandler) CreateProject(c *gin.Context) {
	var input ProjectInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input.Name = strings.TrimSpace(input.Name)

	username, err := h.Jwt.ExtractTokenUsername(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = project.CreateProject(input.Name, username, h.Db, h.ClusterConn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "project " + input.Name + " was created"})
}

func (h *HttpHandler) DeleteProject(c *gin.Context) {
	var input ProjectInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username, err := h.Jwt.ExtractTokenUsername(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = project.DeleteProject(input.Name, username, h.Db, h.ClusterConn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "project " + input.Name + " was deleted"})
}
