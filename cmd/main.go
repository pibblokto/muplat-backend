package main

import (
	"github.com/gin-gonic/gin"
	"github.com/muplat/muplat-backend/controllers"
)

func main() {

	r := gin.Default()

	public := r.Group("/api")

	public.POST("create/ns", controllers.CreateNamespace)

	r.Run(":8080")
}
