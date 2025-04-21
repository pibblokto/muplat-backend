package main

import (
	"github.com/gin-gonic/gin"
	"github.com/muplat/muplat-backend/controllers"
	"github.com/muplat/muplat-backend/middlewares"
	"github.com/muplat/muplat-backend/models"
)

func main() {
	models.ConnectDatabase()
	models.CreateInitUser()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.RedirectTrailingSlash = true

	public := r.Group("/api")
	public.POST("/auth/login", controllers.Login)

	projects := r.Group("/api/project")
	projects.Use(middlewares.JwtAuth())

	projects.POST("", controllers.CreateProject)
	projects.DELETE("", controllers.DeleteProject)

	users := r.Group("/api/auth")
	users.Use(middlewares.JwtAuth())

	users.POST("/user", controllers.AddUser)
	users.DELETE("/user", controllers.DeleteUser)

	r.Run(":8080")
}
