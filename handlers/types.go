package handlers

import (
	"github.com/muplat/muplat-backend/pkg/jwt"
	"github.com/muplat/muplat-backend/pkg/k8s"
	"github.com/muplat/muplat-backend/repositories"
	"github.com/muplat/muplat-backend/services/deployment"
)

type HttpHandler struct {
	Db          *repositories.Database
	ClusterConn *k8s.ClusterConnection
	Jwt         *jwt.JwtConfig
}

func NewHttpHandler(db *repositories.Database, cc *k8s.ClusterConnection, jwt *jwt.JwtConfig) *HttpHandler {
	return &HttpHandler{
		Db:          db,
		ClusterConn: cc,
		Jwt:         jwt,
	}
}

type CreateDeploymentInput struct {
	Name           string                     `json:"name" binding:"required"`
	ProjectName    string                     `json:"projectName" binding:"required"`
	Type           deployment.DeploymentType  `json:"deploymentType" binding:"required"`
	AppConfing     *deployment.AppConfig      `json:"appConfig"`
	PostgresConfig *deployment.PostgresConfig `json:"postgresConfig"`
}

type DeleteDeploymentInput struct {
	Name        string `json:"name" binding:"required"`
	ProjectName string `json:"projectName" binding:"required"`
}
type ReissueCertificateInput struct {
	Name        string `json:"name" binding:"required"`
	ProjectName string `json:"projectName" binding:"required"`
}
type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Admin    *bool  `json:"admin" binding:"required"`
}

type UserDeleteInput struct {
	Username string `json:"username" binding:"required"`
}

type ProjectInput struct {
	Name string `json:"name" binding:"required"`
}

type PatchDeploymentInput struct {
	Type           deployment.DeploymentType       `json:"deploymentType" binding:"required"`
	AppConfing     *deployment.PatchAppConfig      `json:"appConfig"`
	PostgresConfig *deployment.PatchPostgresConfig `json:"postgresConfig"`
}
