package main

import (
	"github.com/gin-gonic/gin"
	"github.com/muplat/muplat-backend/handlers"
	"github.com/muplat/muplat-backend/middlewares"
)

func DeclareRoutes(r *gin.Engine, httpHandler *handlers.HttpHandler) error {
	// public
	public := r.Group("/api")
	public.POST("/auth/login", httpHandler.Login)

	// projects
	projects := r.Group("/api/project")
	projects.Use(middlewares.JwtAuth(httpHandler.Jwt))

	projects.GET("", httpHandler.GetProjects)
	projects.GET("/:project", httpHandler.GetProject)
	projects.POST("", httpHandler.CreateProject)
	projects.DELETE("", httpHandler.DeleteProject)

	// users
	users := r.Group("/api/auth")
	users.Use(middlewares.JwtAuth(httpHandler.Jwt))

	users.GET("/user/:username", httpHandler.GetUser)
	users.GET("/user", httpHandler.GetUsers)
	users.POST("/user", httpHandler.AddUser)
	users.DELETE("/user", httpHandler.DeleteUser)

	// deployments
	deployments := r.Group("/api/deployment")
	deployments.Use(middlewares.JwtAuth(httpHandler.Jwt))

	deployments.GET("/:project/:deployment", httpHandler.GetDeployment)
	deployments.GET("/:project", httpHandler.GetDeployments)
	deployments.POST("/:project/:deployment/reissue", httpHandler.ReissueAppCertificate)
	deployments.POST("", httpHandler.CreateDeployment)
	deployments.DELETE("", httpHandler.DeleteDeployment)
	return nil
}
