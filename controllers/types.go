package controllers

import (
	"github.com/muplat/muplat-backend/pkg/deployments"
	"github.com/muplat/muplat-backend/services"
)

type HttpHandler struct {
	AppService *services.AppService
}

type CreateDeploymentInput struct {
	Name           string                      `json:"name" binding:"required"`
	ProjectName    string                      `json:"projectName" binding:"required"`
	Type           deployments.DeploymentType  `json:"deploymentType" binding:"required"`
	AppConfing     *deployments.AppConfig      `json:"appConfig"`
	PostgresConfig *deployments.PostgresConfig `json:"postgresConfig"`
}

type DeleteDeploymentInput struct {
	Name        string `json:"name" binding:"required"`
	ProjectName string `json:"projectName" binding:"required"`
}
