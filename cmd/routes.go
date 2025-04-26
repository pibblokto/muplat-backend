package main

import (
	"github.com/gin-gonic/gin"
	"github.com/muplat/muplat-backend/controllers"
	"github.com/muplat/muplat-backend/middlewares"
)

func DeclareRoutes(r *gin.Engine, httpHandler *controllers.HttpHandler) error {
	// public
	public := r.Group("/api")
	public.POST("/auth/login", httpHandler.Login)

	// projects
	projects := r.Group("/api/project")
	projects.Use(middlewares.JwtAuth())

	projects.POST("", httpHandler.CreateProject)
	projects.DELETE("", httpHandler.DeleteProject)

	// users
	users := r.Group("/api/auth")
	users.Use(middlewares.JwtAuth())

	users.POST("/user", httpHandler.AddUser)
	users.DELETE("/user", httpHandler.DeleteUser)

	// deployments
	deployments := r.Group("/api/deployment")
	deployments.Use(middlewares.JwtAuth())

	deployments.POST("", httpHandler.CreateDeployment)
	deployments.DELETE("", httpHandler.DeleteDeployment)
	return nil
}
