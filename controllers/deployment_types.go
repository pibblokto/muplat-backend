package controllers

import "github.com/muplat/muplat-backend/pkg/deployments"

type DeploymentInput struct {
	Name        string                     `json:"name" binding:"required"`
	ProjectName string                     `json:"projectName" binding:"required"`
	Type        deployments.DeploymentType `json:"deploymentType" binding:"required"`
	AppConfing  *deployments.AppConfig     `json:"appConfig"`
}
