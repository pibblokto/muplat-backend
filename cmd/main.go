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

	projects := r.Group("/api/project")
	projects.Use(middlewares.JwtAuth())

	projects.POST("", controllers.CreateNamespace)

	r.Run(":8080")
}
