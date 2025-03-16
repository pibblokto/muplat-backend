package main

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/muplat/muplat-backend/controllers"
	"github.com/muplat/muplat-backend/models"
	"github.com/muplat/muplat-backend/middlewares"
)

//"gorm.io/driver/postgres"
//"gorm.io/gorm"

func main() {
	models.ConnectDatabase()

	r := gin.Default()

	public := r.Group("/api")

	public.POST("register", controllers.Register)
	public.POST("login", controllers.Login)

	protected := r.Group("/api/admin")
	protected.Use(middlewares.JwtAuthMiddleware())
	protected.GET("/user", controllers.CurrentUser)
	protected.GET("/secret", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"secret": "yo_mama"})
		fmt.Println(">>> Route handler print")
	})

	r.Run(":8080")
}
